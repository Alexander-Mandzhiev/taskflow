package cors

import (
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var _ contracts.CORSConfig = (*Config)(nil)

type rawConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins" env:"CORS_ALLOWED_ORIGINS"`
	AllowedMethods   []string `mapstructure:"allowed_methods" env:"CORS_ALLOWED_METHODS"`
	AllowedHeaders   []string `mapstructure:"allowed_headers" env:"CORS_ALLOWED_HEADERS"`
	ExposedHeaders   []string `mapstructure:"exposed_headers" env:"CORS_EXPOSED_HEADERS"`
	AllowCredentials bool     `mapstructure:"allow_credentials" env:"CORS_ALLOW_CREDENTIALS"`
	MaxAge           int      `mapstructure:"max_age" env:"CORS_MAX_AGE"`
}

// Config — конфиг CORS (политика доступа для HTTP).
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		AllowedOrigins:   nil,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Accept", "X-Requested-With", "Origin"},
		ExposedHeaders:   nil,
		AllowCredentials: false,
		MaxAge:           0,
	}
}

func (c *Config) AllowedOrigins() []string { return c.raw.AllowedOrigins }
func (c *Config) AllowedMethods() []string { return c.raw.AllowedMethods }
func (c *Config) AllowedHeaders() []string { return c.raw.AllowedHeaders }
func (c *Config) ExposedHeaders() []string { return c.raw.ExposedHeaders }
func (c *Config) AllowCredentials() bool   { return c.raw.AllowCredentials }
func (c *Config) MaxAge() int              { return c.raw.MaxAge }
