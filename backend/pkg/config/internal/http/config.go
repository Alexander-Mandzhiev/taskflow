package http

import (
	"fmt"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var _ contracts.HTTPConfig = (*Config)(nil)
var _ contracts.Validatable = (*Config)(nil)

type rawConfig struct {
	Address           string        `mapstructure:"address" env:"HTTP_ADDRESS"`
	Timeout           time.Duration `mapstructure:"timeout" env:"HTTP_TIMEOUT"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout" env:"HTTP_READ_HEADER_TIMEOUT"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout" env:"HTTP_READ_TIMEOUT"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout" env:"HTTP_IDLE_TIMEOUT"`
	MaxHeaderBytes    int           `mapstructure:"max_header_bytes" env:"HTTP_MAX_HEADER_BYTES"`
	ShutdownTimeout   time.Duration `mapstructure:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT"`
}

// Config — конфиг HTTP-сервера.
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		Address:           ":8080",
		Timeout:           30 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:  1 << 20, // 1 MiB
		ShutdownTimeout: 10 * time.Second,
	}
}

func (c *Config) Address() string              { return c.raw.Address }
func (c *Config) Timeout() time.Duration      { return c.raw.Timeout }
func (c *Config) ReadHeaderTimeout() time.Duration { return c.raw.ReadHeaderTimeout }
func (c *Config) ReadTimeout() time.Duration  { return c.raw.ReadTimeout }
func (c *Config) WriteTimeout() time.Duration { return c.raw.WriteTimeout }
func (c *Config) IdleTimeout() time.Duration  { return c.raw.IdleTimeout }
func (c *Config) MaxHeaderBytes() int         { return c.raw.MaxHeaderBytes }
func (c *Config) ShutdownTimeout() time.Duration { return c.raw.ShutdownTimeout }

// Validate проверяет корректность настроек HTTP.
func (c *Config) Validate() error {
	if c.raw.Address == "" {
		return fmt.Errorf("address must be non-empty")
	}
	if c.raw.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if c.raw.MaxHeaderBytes <= 0 {
		return fmt.Errorf("max_header_bytes must be positive")
	}
	return nil
}
