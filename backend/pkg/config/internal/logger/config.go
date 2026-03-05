package logger

import (
	"fmt"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var _ contracts.LoggerConfig = (*Config)(nil)
var _ contracts.Validatable = (*Config)(nil)

type rawConfig struct {
	Level                string        `mapstructure:"level" env:"LOGGER_LEVEL"`
	AsJSON               bool          `mapstructure:"as_json" env:"LOGGER_AS_JSON"`
	Name                 string        `mapstructure:"name" env:"LOGGER_NAME"`
	Environment          string        `mapstructure:"environment" env:"LOGGER_ENVIRONMENT"`
	OTLPEnable           bool          `mapstructure:"otlp_enable" env:"LOGGER_OTLP_ENABLE"`
	OTLPEndpoint         string        `mapstructure:"otlp_endpoint" env:"LOGGER_OTLP_ENDPOINT"`
	OTLPShutdownTimeout  time.Duration `mapstructure:"otlp_shutdown_timeout" env:"LOGGER_OTLP_SHUTDOWN_TIMEOUT"`
}

// Config — конфиг модуля logger.
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		Level:                "info",
		AsJSON:               true,
		Name:                 "mkk",
		Environment:          "development",
		OTLPEnable:           false,
		OTLPEndpoint:         "localhost:4317",
		OTLPShutdownTimeout:  30 * time.Second,
	}
}

func (c *Config) Level() string                { return c.raw.Level }
func (c *Config) AsJSON() bool                 { return c.raw.AsJSON }
func (c *Config) Name() string                 { return c.raw.Name }
func (c *Config) Environment() string          { return c.raw.Environment }
func (c *Config) OTLPEnable() bool             { return c.raw.OTLPEnable }
func (c *Config) OTLPEndpoint() string         { return c.raw.OTLPEndpoint }
func (c *Config) OTLPShutdownTimeout() time.Duration { return c.raw.OTLPShutdownTimeout }

// Validate проверяет корректность настроек логгера.
func (c *Config) Validate() error {
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.raw.Level] {
		return fmt.Errorf("level must be one of: debug, info, warn, error")
	}
	return nil
}
