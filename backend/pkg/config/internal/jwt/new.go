package jwt

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/helpers"
)

// New создаёт конфиг JWT по стратегии: Defaults → YAML → ENV.
func New() (contracts.JWTConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("jwt"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal jwt YAML: %w", err)
		}
	}
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse jwt ENV: %w", err)
	}
	return cfg, nil
}
