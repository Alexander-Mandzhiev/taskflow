package http

import (
	"fmt"
	"time"

	"mkk/pkg/config/contracts"
)

var _ contracts.HTTPConfig = (*Config)(nil)
var _ contracts.Validatable = (*Config)(nil)

type rawConfig struct {
	Address string        `mapstructure:"address" env:"HTTP_ADDRESS"`
	Timeout time.Duration `mapstructure:"timeout" env:"HTTP_TIMEOUT"`
}

// Config — конфиг HTTP-сервера.
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		Address: ":8080",
		Timeout: 30 * time.Second,
	}
}

func (c *Config) Address() string        { return c.raw.Address }
func (c *Config) Timeout() time.Duration { return c.raw.Timeout }

// Validate проверяет корректность настроек HTTP.
func (c *Config) Validate() error {
	if c.raw.Address == "" {
		return fmt.Errorf("address must be non-empty")
	}
	if c.raw.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	return nil
}
