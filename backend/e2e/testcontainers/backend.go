//go:build integration

package testcontainers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Backend — контейнер приложения (бэкенд API).
// Подключается к MySQL и Redis через Docker-сеть по алиасам контейнеров.
type Backend struct {
	ctr testcontainers.Container
	url string
}

// BackendOpts — опции запуска контейнера бэкенда (пути к Dockerfile и контексту сборки).
type BackendOpts struct {
	BuildContext   string // пусто — defaultBuildContext
	DockerfilePath string // пусто — defaultDockerfile
}

// NewBackend собирает образ из Dockerfile и поднимает контейнер в указанной Docker-сети.
// MySQL и Redis доступны внутри сети по алиасам MySQLContainerAlias и RedisContainerAlias.
func NewBackend(ctx context.Context, networkName string, opts BackendOpts) (*Backend, error) {
	buildContext := opts.BuildContext
	if buildContext == "" {
		projectRoot, err := findProjectRoot()
		if err != nil {
			return nil, fmt.Errorf("project root: %w", err)
		}
		buildContext = projectRoot
	}
	dockerfilePath := opts.DockerfilePath
	if dockerfilePath == "" {
		dockerfilePath = defaultDockerfile
	}

	// Бэкенд подключается к MySQL и Redis по именам контейнеров внутри Docker-сети,
	// поэтому не нужны host.docker.internal и динамические порты хоста.
	appEnv := map[string]string{
		"CONFIG_PATH":     "/app/config/test.yaml",
		"MIGRATIONS_DIR":  "db/migration",
		"MIGRATIONS_AUTO": "true",
		"MYSQL_HOST":      MySQLContainerAlias,
		"MYSQL_PORT":      MySQLInternalPort,
		"MYSQL_USER":      MySQLUser,
		"MYSQL_PASSWORD":  MySQLPassword,
		"MYSQL_DATABASE":  MySQLDatabase,
		"REDIS_ADDR":      RedisContainerAlias + ":" + RedisInternalPort,
		"REDIS_PASSWORD":  RedisPassword,
		"HTTP_ADDRESS":    "0.0.0.0:8080",
	}

	req := testcontainers.GenericContainerRequest{
		ProviderType: testcontainers.ProviderDocker,
		Started:      true,
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    buildContext,
				Dockerfile: dockerfilePath,
				// Стабильное имя образа + KeepImage: Docker использует layer cache при повторных запусках.
				Repo:      BackendImageRepo,
				Tag:       BackendImageTag,
				KeepImage: true,
			},
			Env:          appEnv,
			ExposedPorts: []string{BackendHTTPPort + "/tcp"},
			Networks:     []string{networkName},
			WaitingFor: wait.ForHTTP("/health").
				WithPort(BackendHTTPPort + "/tcp").
				WithStatusCodeMatcher(func(status int) bool { return status == http.StatusOK }).
				WithStartupTimeout(backendStartupLimit),
			HostConfigModifier: func(config *container.HostConfig) {
				config.LogConfig = container.LogConfig{
					Type:   "json-file",
					Config: map[string]string{},
				}
			},
		},
	}

	ctr, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("backend container: %w", err)
	}

	port, portErr := ctr.MappedPort(ctx, nat.Port(BackendHTTPPort+"/tcp"))
	host, hostErr := ctr.Host(ctx)
	if portErr != nil || hostErr != nil {
		dumpBackendDiagnostics(ctx, ctr, portErr, hostErr, appEnv)
		_ = ctr.Terminate(ctx)
		if portErr != nil {
			return nil, fmt.Errorf("backend port: %w", portErr)
		}
		return nil, fmt.Errorf("backend host: %w", hostErr)
	}

	url := "http://" + net.JoinHostPort(host, port.Port())

	dumpBackendLogs(ctx, ctr, "post-start")
	go func() {
		time.Sleep(3 * time.Second)
		dumpBackendLogs(context.Background(), ctr, "debug-after-3s")
	}()

	return &Backend{ctr: ctr, url: url}, nil
}

// dumpBackendLogs выводит последние логи контейнера (для отладки падения бэкенда).
func dumpBackendLogs(ctx context.Context, ctr testcontainers.Container, reason string) {
	logsRC, err := ctr.Logs(ctx)
	if err != nil {
		log.Printf("[testcontainers] %s: не удалось получить логи контейнера: %v", reason, err)
		return
	}
	defer func() { _ = logsRC.Close() }()
	const maxLogBytes = 16384
	b, readErr := io.ReadAll(io.LimitReader(logsRC, maxLogBytes))
	if readErr != nil {
		log.Printf("[testcontainers] %s: ошибка чтения логов: %v", reason, readErr)
		return
	}
	if len(b) > 0 {
		log.Printf("[testcontainers] %s: логи бэкенда:\n%s", reason, string(b))
	} else {
		log.Printf("[testcontainers] %s: логов нет (контейнер мог завершиться до вывода)", reason)
	}
}

type containerIDGetter interface {
	GetContainerID() string
}

// dumpBackendDiagnostics выводит диагностику при ошибке старта.
func dumpBackendDiagnostics(ctx context.Context, ctr testcontainers.Container, portErr, hostErr error, appEnv map[string]string) {
	log.Println("[testcontainers] === диагностика падения контейнера бэкенда ===")
	if portErr != nil {
		log.Printf("[testcontainers] ошибка порта: %v", portErr)
	}
	if hostErr != nil {
		log.Printf("[testcontainers] ошибка хоста: %v", hostErr)
	}
	dumpBackendLogs(ctx, ctr, "логи при падении")
	if getter, ok := ctr.(containerIDGetter); ok {
		if id := getter.GetContainerID(); id != "" {
			if out, err := runDockerInspect(ctx, id); err != nil {
				log.Printf("[testcontainers] docker inspect: %v", err)
			} else if out != "" {
				log.Printf("[testcontainers] docker inspect (состояние/exit code):\n%s", out)
			}
		}
	}
	log.Println("[testcontainers] переменные окружения контейнера:")
	for k, v := range appEnv {
		log.Printf("[testcontainers]   %s=%s", k, v)
	}
	log.Println("[testcontainers] === конец диагностики ===")
}

func runDockerInspect(ctx context.Context, containerID string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "inspect", "--format",
		"Status: {{.State.Status}}, ExitCode: {{.State.ExitCode}}, Error: {{.State.Error}}", containerID)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Terminate останавливает контейнер.
func (b *Backend) Terminate(ctx context.Context) error {
	if b.ctr == nil {
		return nil
	}
	return b.ctr.Terminate(ctx)
}

// URL возвращает base URL API для apiclient (например http://localhost:8080).
func (b *Backend) URL() string { return b.url }

// findProjectRoot возвращает корень репо (родитель backend/), где лежит deploy/docker/backend/.
func findProjectRoot() (string, error) {
	wd, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			parent := filepath.Dir(dir)
			if parent == dir {
				return "", fmt.Errorf("cannot find project root from %s", wd)
			}
			return parent, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", wd)
		}
		dir = parent
	}
}
