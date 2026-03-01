package cache

import (
	"time"
)

// Option функциональная опция для конфигурации Redis
type Option func(*Config)

// Config содержит все настройки подключения к Redis (один узел)
type Config struct {
	addr            string
	password        string
	poolSize        int
	minIdleConns    int
	poolTimeout     time.Duration
	connMaxIdleTime time.Duration
	dialTimeout     time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration
}

// defaultConfig возвращает конфигурацию по умолчанию (без YAML/ENV приложение подключается к локальному Redis).
func defaultConfig() *Config {
	return &Config{
		addr:            "localhost:6379",
		password:        "",
		poolSize:        10,
		minIdleConns:    2,
		poolTimeout:     4 * time.Second,
		connMaxIdleTime: 5 * time.Minute,
		dialTimeout:     5 * time.Second,
		readTimeout:     3 * time.Second,
		writeTimeout:    3 * time.Second,
	}
}

// WithAddr задаёт адрес Redis (например "localhost:6379"). Пустой адрес — кеш отключён.
func WithAddr(addr string) Option {
	return func(c *Config) { c.addr = addr }
}

// WithPassword задаёт пароль для подключения к Redis. Не экспортируется в структуре — только в redis.Options.
func WithPassword(p string) Option {
	return func(c *Config) { c.password = p }
}

// WithPoolSize задаёт размер пула соединений.
func WithPoolSize(n int) Option {
	return func(c *Config) { c.poolSize = n }
}

// WithMinIdleConns задаёт минимальное количество idle-соединений в пуле.
func WithMinIdleConns(n int) Option {
	return func(c *Config) { c.minIdleConns = n }
}

// WithPoolTimeout задаёт таймаут ожидания соединения из пула.
func WithPoolTimeout(d time.Duration) Option {
	return func(c *Config) { c.poolTimeout = d }
}

// WithConnMaxIdleTime задаёт максимальное время жизни idle-соединения.
func WithConnMaxIdleTime(d time.Duration) Option {
	return func(c *Config) { c.connMaxIdleTime = d }
}

// WithDialTimeout задаёт таймаут установки соединения.
func WithDialTimeout(d time.Duration) Option {
	return func(c *Config) { c.dialTimeout = d }
}

// WithReadTimeout задаёт таймаут на чтение.
func WithReadTimeout(d time.Duration) Option {
	return func(c *Config) { c.readTimeout = d }
}

// WithWriteTimeout задаёт таймаут на запись.
func WithWriteTimeout(d time.Duration) Option {
	return func(c *Config) { c.writeTimeout = d }
}
