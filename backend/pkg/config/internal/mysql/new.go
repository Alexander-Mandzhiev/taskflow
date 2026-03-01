package mysql

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"mkk/pkg/config/contracts"
	"mkk/pkg/config/helpers"
)

// New создаёт конфиг mysql по стратегии: Defaults → YAML → ENV.
func New() (contracts.MySQLConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("mysql"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal mysql YAML: %w", err)
		}
	}

	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse mysql ENV: %w", err)
	}
	cfg.dsn = cfg.buildDSN()
	return cfg, nil
}
