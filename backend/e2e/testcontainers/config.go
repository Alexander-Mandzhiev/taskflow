package testcontainers

import (
	"context"
	"fmt"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

// LoadTestConfig выставляет адреса контейнеров в env (ApplyTestEnv), затем вызывает config.Load(ctx).
// User, password, database и redis password берутся из config/test.yaml (pkg/config).
func LoadTestConfig(ctx context.Context, env *TestEnv) (contracts.Provider, error) {
	ApplyTestEnv(env)
	cfg, err := config.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("config.Load: %w", err)
	}
	return cfg, nil
}
