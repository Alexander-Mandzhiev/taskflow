package redis

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"mkk/pkg/config/contracts"
	"mkk/pkg/config/helpers"
)

// New создаёт конфиг redis по стратегии: Defaults → YAML → ENV.
func New() (contracts.RedisConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("redis"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal redis YAML: %w", err)
		}
	}

	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse redis ENV: %w", err)
	}
	cfg.poolConfig = &Pool{raw: cfg.raw.Pool}
	return cfg, nil
}
