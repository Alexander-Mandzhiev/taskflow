package connectingpool

import (
	"time"
)

// Option — функциональная опция для конфигурации пула соединений.
type Option func(*Config)

// Config — настройки пула соединений к БД.
type Config struct {
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

func defaultConfig() *Config {
	return &Config{
		maxOpenConns:    25,
		maxIdleConns:    5,
		connMaxLifetime: 5 * time.Minute,
		connMaxIdleTime: 3 * time.Minute,
	}
}

// WithMaxOpenConns задаёт максимальное количество открытых соединений (0 = без лимита).
func WithMaxOpenConns(n int) Option {
	return func(c *Config) { c.maxOpenConns = n }
}

// WithMaxIdleConns задаёт максимальное количество idle-соединений в пуле.
func WithMaxIdleConns(n int) Option {
	return func(c *Config) { c.maxIdleConns = n }
}

// WithConnMaxLifetime задаёт максимальное время жизни соединения.
func WithConnMaxLifetime(d time.Duration) Option {
	return func(c *Config) { c.connMaxLifetime = d }
}

// WithConnMaxIdleTime задаёт максимальное время простоя соединения перед закрытием.
func WithConnMaxIdleTime(d time.Duration) Option {
	return func(c *Config) { c.connMaxIdleTime = d }
}
