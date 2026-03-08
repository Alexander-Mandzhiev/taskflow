//go:build integration

package testcontainers

import (
	"context"
	"fmt"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

// Redis — контейнер Redis и адрес подключения (host:port).
type Redis struct {
	ctr  *tcredis.RedisContainer
	addr string
}

// NewRedis поднимает контейнер Redis и подключает его к указанной Docker-сети.
func NewRedis(ctx context.Context, networkName string) (*Redis, error) {
	ctr, err := tcredis.Run(ctx, RedisImage,
		tcredis.WithSnapshotting(10, 1),
		testcontainers.WithCmdArgs("--requirepass", RedisPassword),
		withNetworkAlias(networkName, RedisContainerAlias),
	)
	if err != nil {
		return nil, fmt.Errorf("redis.Run: %w", err)
	}

	connStr, err := ctr.ConnectionString(ctx)
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, fmt.Errorf("redis ConnectionString: %w", err)
	}

	addr := connStr
	if len(addr) > 8 && addr[:8] == "redis://" {
		addr = addr[8:]
	}
	addr = strings.TrimSpace(addr)

	return &Redis{ctr: ctr, addr: addr}, nil
}

// Terminate останавливает контейнер.
func (r *Redis) Terminate(ctx context.Context) error {
	if r.ctr == nil {
		return nil
	}
	return r.ctr.Terminate(ctx)
}

// Addr возвращает адрес host:port для подключения.
func (r *Redis) Addr() string { return r.addr }
