package jwt

import (
	"fmt"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var (
	_ contracts.JWTConfig   = (*Config)(nil)
	_ contracts.Validatable = (*Config)(nil)
)

type rawConfig struct {
	AccessSecret           string        `mapstructure:"access_secret" env:"JWT_ACCESS_SECRET"`   //nolint:gosec
	RefreshSecret          string        `mapstructure:"refresh_secret" env:"JWT_REFRESH_SECRET"` //nolint:gosec // G117: конфиг/env, не логируем
	AccessTTL              time.Duration `mapstructure:"access_ttl" env:"JWT_ACCESS_TTL"`
	RefreshTTL             time.Duration `mapstructure:"refresh_ttl" env:"JWT_REFRESH_TTL"`
	AccessTokenCookieName  string        `mapstructure:"access_token_cookie_name" env:"JWT_ACCESS_TOKEN_COOKIE_NAME"`
	RefreshTokenCookieName string        `mapstructure:"refresh_token_cookie_name" env:"JWT_REFRESH_TOKEN_COOKIE_NAME"`
}

// Config — конфиг JWT.
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		AccessSecret:           "",
		RefreshSecret:          "",
		AccessTTL:              15 * time.Minute,
		RefreshTTL:             7 * 24 * time.Hour, // 7 дней
		AccessTokenCookieName:  "access_token",
		RefreshTokenCookieName: "refresh_token",
	}
}

func (c *Config) AccessSecret() string { return c.raw.AccessSecret }
func (c *Config) RefreshSecret() string {
	if c.raw.RefreshSecret != "" {
		return c.raw.RefreshSecret
	}
	return c.raw.AccessSecret
}
func (c *Config) AccessTTL() time.Duration  { return c.raw.AccessTTL }
func (c *Config) RefreshTTL() time.Duration { return c.raw.RefreshTTL }
func (c *Config) AccessTokenCookieName() string {
	if c.raw.AccessTokenCookieName != "" {
		return c.raw.AccessTokenCookieName
	}
	return "access_token"
}

func (c *Config) RefreshTokenCookieName() string {
	if c.raw.RefreshTokenCookieName != "" {
		return c.raw.RefreshTokenCookieName
	}
	return "refresh_token"
}

// Validate проверяет корректность настроек JWT.
func (c *Config) Validate() error {
	if c.raw.AccessSecret == "" {
		return fmt.Errorf("jwt access_secret cannot be empty")
	}
	if c.raw.AccessTTL <= 0 {
		return fmt.Errorf("jwt access_ttl must be positive")
	}
	if c.raw.RefreshTTL <= 0 {
		return fmt.Errorf("jwt refresh_ttl must be positive")
	}
	if c.AccessTokenCookieName() == "" {
		return fmt.Errorf("jwt access_token_cookie_name cannot be empty")
	}
	if c.RefreshTokenCookieName() == "" {
		return fmt.Errorf("jwt refresh_token_cookie_name cannot be empty")
	}
	return nil
}
