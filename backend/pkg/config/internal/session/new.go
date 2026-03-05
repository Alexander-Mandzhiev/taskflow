package session

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/helpers"
)

// New создаёт конфиг session по стратегии: Defaults → YAML → ENV.
func New() (contracts.SessionConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("session"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal session YAML: %w", err)
		}
	}
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse session ENV: %w", err)
	}
	return cfg, nil
}
