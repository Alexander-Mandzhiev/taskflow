package session

import (
	"fmt"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var (
	_ contracts.SessionConfig = (*Config)(nil)
	_ contracts.Validatable   = (*Config)(nil)
)

type rawConfig struct {
	TTL          time.Duration `mapstructure:"ttl" env:"SESSION_TTL"`
	IsSecure     bool          `mapstructure:"is_secure" env:"SESSION_IS_SECURE"`
	CookieDomain string        `mapstructure:"cookie_domain" env:"SESSION_COOKIE_DOMAIN"`
}

// Config — конфиг модуля session.
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		TTL:          24 * time.Hour,
		IsSecure:     false,
		CookieDomain: "",
	}
}

func (c *Config) TTL() time.Duration   { return c.raw.TTL }
func (c *Config) IsSecure() bool       { return c.raw.IsSecure }
func (c *Config) CookieDomain() string { return c.raw.CookieDomain }

// Validate проверяет корректность настроек сессии.
func (c *Config) Validate() error {
	if c.raw.TTL <= 0 {
		return fmt.Errorf("ttl must be positive")
	}
	return nil
}
