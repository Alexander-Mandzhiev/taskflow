package app

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/helpers"
)

// New создаёт конфиг app по стратегии: Defaults → YAML → ENV.
func New() (contracts.AppConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("app"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal app YAML: %w", err)
		}
	}
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse app ENV: %w", err)
	}
	return cfg, nil
}
