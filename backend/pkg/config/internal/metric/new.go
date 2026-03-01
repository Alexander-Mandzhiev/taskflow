package metric

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"mkk/pkg/config/contracts"
	"mkk/pkg/config/helpers"
)

// New создаёт конфиг metric по стратегии: Defaults → YAML → ENV.
func New() (contracts.MetricConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("metric"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal metric YAML: %w", err)
		}
	}
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse metric ENV: %w", err)
	}
	return cfg, nil
}
