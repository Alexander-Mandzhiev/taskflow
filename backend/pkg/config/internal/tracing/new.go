package tracing

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"mkk/pkg/config/contracts"
	"mkk/pkg/config/helpers"
)

// New создаёт конфиг tracing по стратегии: Defaults → YAML → ENV.
func New() (contracts.TracingConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("tracing"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal tracing YAML: %w", err)
		}
	}
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse tracing ENV: %w", err)
	}
	return cfg, nil
}
