package cors

import (
	"fmt"

	"github.com/caarlos0/env/v11"

	"mkk/pkg/config/contracts"
	"mkk/pkg/config/helpers"
)

// New создаёт конфиг CORS по стратегии: Defaults → YAML → ENV.
func New() (contracts.CORSConfig, error) {
	cfg := &Config{raw: defaultConfig()}
	if section := helpers.GetSection("cors"); section != nil {
		if err := section.Unmarshal(&cfg.raw); err != nil {
			return nil, fmt.Errorf("unmarshal cors YAML: %w", err)
		}
	}
	if err := env.Parse(&cfg.raw); err != nil {
		return nil, fmt.Errorf("parse cors ENV: %w", err)
	}
	return cfg, nil
}
