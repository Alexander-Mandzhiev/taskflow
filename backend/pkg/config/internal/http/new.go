package http

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"mkk/pkg/config/contracts"
	"mkk/pkg/config/helpers"
)

// New создаёт конфиг http по стратегии: Defaults → YAML → ENV.
func New() (contracts.HTTPConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("http"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal http YAML: %w", err)
		}
	}
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse http ENV: %w", err)
	}
	return cfg, nil
}
