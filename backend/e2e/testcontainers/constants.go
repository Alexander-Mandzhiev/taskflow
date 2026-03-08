//go:build integration

package testcontainers

import (
	"time"

	tc "github.com/testcontainers/testcontainers-go"
)

// Креды и образы согласованы с backend/config/test.yaml (единственный источник правды для переменных при config.Load).
const (
	MySQLImage    = "mysql:8.4"
	MySQLDatabase = "testdb"
	MySQLUser     = "testuser"
	MySQLPassword = "testpass"

	RedisImage    = "redis:7.2.5-alpine"
	RedisPassword = "test_redis_secret"
)

// Backend — порт HTTP API (должен совпадать с http.address в config/test.yaml) и таймаут старта контейнера.
const (
	BackendHTTPPort     = "8080"
	backendStartupLimit = 5 * time.Minute
	defaultBuildContext = ".."
	defaultDockerfile   = "deploy/docker/backend/Dockerfile"

	// Имя образа для кэширования между запусками (KeepImage: true).
	BackendImageRepo = "taskflow-backend-test"
	BackendImageTag  = "latest"
)

// Docker-сеть и алиасы контейнеров внутри неё.
const (
	NetworkName         = "taskflow-e2e"
	MySQLContainerAlias = "mysql"
	RedisContainerAlias = "redis"
	MySQLInternalPort   = "3306"
	RedisInternalPort   = "6379"
)

// withNetworkAlias возвращает опцию, добавляющую контейнер в Docker-сеть с заданным alias.
// Используется при запуске модулей mysql и redis.
func withNetworkAlias(networkName, alias string) tc.ContainerCustomizer {
	return tc.CustomizeRequestOption(func(req *tc.GenericContainerRequest) error {
		req.Networks = append(req.Networks, networkName)
		if req.NetworkAliases == nil {
			req.NetworkAliases = make(map[string][]string)
		}
		req.NetworkAliases[networkName] = []string{alias}
		return nil
	})
}
