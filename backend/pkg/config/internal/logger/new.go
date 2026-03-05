package logger

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/helpers"
)

// New создаёт конфиг logger по стратегии: Defaults → YAML → ENV.
func New() (contracts.LoggerConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("logger"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal logger YAML: %w", err)
		}
	}
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse logger ENV: %w", err)
	}
	return cfg, nil
}
