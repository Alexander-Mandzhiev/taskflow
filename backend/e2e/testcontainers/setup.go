//go:build integration

package testcontainers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
)

// TestEnvironment — полное тестовое окружение (сеть, MySQL, Redis, Backend).
// Данные для подключения в TestEnv, остановка — через Cleanup(ctx).
type TestEnvironment struct {
	TestEnv
	dockerNetwork testcontainers.Network
	mysql         *MySQL
	redis         *Redis
	backend       *Backend
}

// Setup поднимает полный стек в порядке зависимостей:
//
//  1. Docker network        — изолированная сеть для всех контейнеров
//  2. MySQL (mysql.go)      — БД, подключена к сети с алиасом "mysql"
//  3. Redis (redis.go)      — кеш/сессии, подключён к сети с алиасом "redis"
//  4. Ping deps             — проверка доступности с хоста через маппированные порты
//  5. Backend (backend.go)  — приложение, подключается к MySQL/Redis через имена контейнеров в сети
//
// Остановка — Cleanup(ctx): Backend → Redis → MySQL → Network.
func Setup(ctx context.Context) (*TestEnvironment, error) {
	net, err := testcontainers.GenericNetwork(ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           NetworkName,
			Driver:         "bridge",
			CheckDuplicate: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("create docker network %q: %w", NetworkName, err)
	}

	mysq, err := NewMySQL(ctx, NetworkName)
	if err != nil {
		_ = net.Remove(ctx)
		return nil, err
	}

	redi, err := NewRedis(ctx, NetworkName)
	if err != nil {
		_ = mysq.Terminate(ctx)
		_ = net.Remove(ctx)
		return nil, err
	}

	if err := pingDeps(ctx, mysq.DSN(), redi.Addr()); err != nil {
		_ = redi.Terminate(ctx)
		_ = mysq.Terminate(ctx)
		_ = net.Remove(ctx)
		return nil, fmt.Errorf("deps not ready: %w", err)
	}

	app, err := NewBackend(ctx, NetworkName, BackendOpts{})
	if err != nil {
		_ = redi.Terminate(ctx)
		_ = mysq.Terminate(ctx)
		_ = net.Remove(ctx)
		return nil, err
	}

	env := &TestEnvironment{
		TestEnv: TestEnv{
			MySQLHost:  mysq.Host(),
			MySQLPort:  mysq.Port(),
			DSN:        mysq.DSN(),
			RedisAddr:  redi.Addr(),
			BackendURL: app.URL(),
		},
		dockerNetwork: net,
		mysql:         mysq,
		redis:         redi,
		backend:       app,
	}
	return env, nil
}

// Cleanup останавливает контейнеры и удаляет сеть: Backend → Redis → MySQL → Network.
func (e *TestEnvironment) Cleanup(ctx context.Context) {
	if e.backend != nil {
		if err := e.backend.Terminate(ctx); err != nil {
			log.Printf("[testcontainers] не удалось остановить контейнер приложения: %v", err)
		} else {
			log.Println("[testcontainers] 🛑 Контейнер приложения остановлен")
		}
	}
	if e.redis != nil {
		if err := e.redis.Terminate(ctx); err != nil {
			log.Printf("[testcontainers] не удалось остановить контейнер Redis: %v", err)
		} else {
			log.Println("[testcontainers] 🛑 Контейнер Redis остановлен")
		}
	}
	if e.mysql != nil {
		if err := e.mysql.Terminate(ctx); err != nil {
			log.Printf("[testcontainers] не удалось остановить контейнер MySQL: %v", err)
		} else {
			log.Println("[testcontainers] 🛑 Контейнер MySQL остановлен")
		}
	}
	if e.dockerNetwork != nil {
		if err := e.dockerNetwork.Remove(ctx); err != nil {
			log.Printf("[testcontainers] не удалось удалить Docker-сеть: %v", err)
		} else {
			log.Println("[testcontainers] 🛑 Docker-сеть удалена")
		}
	}
}

// pingDeps проверяет, что MySQL и Redis принимают соединения (пинг с хоста через маппированные порты).
func pingDeps(ctx context.Context, dsn, redisAddr string) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("mysql open: %w", err)
	}
	defer db.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("mysql ping: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: RedisPassword,
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	return nil
}
