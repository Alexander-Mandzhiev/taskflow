package app

import (
	"fmt"

	"mkk/pkg/config/contracts"
)

var _ contracts.AppConfig = (*Config)(nil)
var _ contracts.Validatable = (*Config)(nil)

type rawConfig struct {
	Name        string `mapstructure:"name" env:"APP_NAME"`
	Environment string `mapstructure:"environment" env:"APP_ENVIRONMENT"`
	Version     string `mapstructure:"version" env:"APP_VERSION"`
}

// Config — конфиг модуля app.
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		Name:        "mkk",
		Environment: "development",
		Version:     "1.0.0",
	}
}

func (c *Config) Name() string        { return c.raw.Name }
func (c *Config) Environment() string { return c.raw.Environment }
func (c *Config) Version() string     { return c.raw.Version }

// Validate проверяет корректность настроек приложения.
func (c *Config) Validate() error {
	if c.raw.Name == "" {
		return fmt.Errorf("name must be non-empty")
	}
	return nil
}
