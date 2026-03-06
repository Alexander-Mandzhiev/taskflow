package app

import (
	"fmt"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var (
	_ contracts.AppConfig   = (*Config)(nil)
	_ contracts.Validatable = (*Config)(nil)
)

type rawConfig struct {
	Name          string `mapstructure:"name" env:"APP_NAME"`
	Environment   string `mapstructure:"environment" env:"APP_ENVIRONMENT"`
	Version       string `mapstructure:"version" env:"APP_VERSION"`
	MigrationsDir string `mapstructure:"migrations_dir" env:"MIGRATIONS_DIR"`
}

// Config — конфиг модуля app.
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		Name:          "mkk",
		Environment:   "development",
		Version:       "1.0.0",
		MigrationsDir: "db/migration",
	}
}

func (c *Config) Name() string          { return c.raw.Name }
func (c *Config) Environment() string   { return c.raw.Environment }
func (c *Config) Version() string       { return c.raw.Version }
func (c *Config) MigrationsDir() string { return c.raw.MigrationsDir }

// Validate проверяет корректность настроек приложения.
func (c *Config) Validate() error {
	if c.raw.Name == "" {
		return fmt.Errorf("name must be non-empty")
	}
	return nil
}
