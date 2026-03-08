//go:build integration

package testcontainers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// MySQL — контейнер MySQL и параметры подключения (host, port, DSN).
type MySQL struct {
	ctr  *mysql.MySQLContainer
	host string
	port int
	dsn  string
}

// NewMySQL поднимает контейнер MySQL и подключает его к указанной Docker-сети.
func NewMySQL(ctx context.Context, networkName string) (*MySQL, error) {
	ctr, err := mysql.Run(ctx, MySQLImage,
		mysql.WithDatabase(MySQLDatabase),
		mysql.WithUsername(MySQLUser),
		mysql.WithPassword(MySQLPassword),
		withNetworkAlias(networkName, MySQLContainerAlias),
	)
	if err != nil {
		return nil, fmt.Errorf("mysql.Run: %w", err)
	}

	dsn, err := ctr.ConnectionString(ctx)
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, fmt.Errorf("mysql ConnectionString: %w", err)
	}

	host, err := ctr.Host(ctx)
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, fmt.Errorf("mysql Host: %w", err)
	}

	portStr, err := ctr.MappedPort(ctx, "3306")
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, fmt.Errorf("mysql MappedPort: %w", err)
	}

	port, err := strconv.Atoi(portStr.Port())
	if err != nil {
		_ = ctr.Terminate(ctx)
		return nil, fmt.Errorf("mysql port parse %q: %w", portStr.Port(), err)
	}
	if port <= 0 {
		_ = ctr.Terminate(ctx)
		return nil, fmt.Errorf("mysql port invalid: %d", port)
	}

	return &MySQL{ctr: ctr, host: host, port: port, dsn: dsn}, nil
}

// Terminate останавливает контейнер.
func (m *MySQL) Terminate(ctx context.Context) error {
	if m.ctr == nil {
		return nil
	}
	return m.ctr.Terminate(ctx)
}

// Host возвращает хост для подключения (с хоста теста — localhost).
func (m *MySQL) Host() string { return m.host }

// Port возвращает маппированный порт.
func (m *MySQL) Port() int { return m.port }

// DSN возвращает connection string для MySQL.
func (m *MySQL) DSN() string { return m.dsn }
