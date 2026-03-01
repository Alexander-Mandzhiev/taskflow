package metric

import (
	"time"
)

// Option функциональная опция для конфигурации метрик
type Option func(*Config)

// Config содержит все настройки метрик
type Config struct {
	name        string
	environment string
	version     string
	enable      bool
	namespace   string
	appName     string

	// Настройки OpenTelemetry
	endpoint        string
	timeout         time.Duration
	exportInterval  time.Duration
	shutdownTimeout time.Duration
}

// defaultConfig возвращает конфигурацию по умолчанию
func defaultConfig() *Config {
	return &Config{
		name:            "service",
		environment:     "development",
		version:         "1.0.0",
		enable:          false,
		namespace:       "my_space",
		appName:         "my_app",
		endpoint:        "localhost:4317",
		timeout:         30 * time.Second,
		exportInterval:  5 * time.Second,
		shutdownTimeout: 30 * time.Second,
	}
}

// WithName устанавливает имя сервиса
func WithName(name string) Option {
	return func(c *Config) {
		c.name = name
	}
}

// WithEnvironment устанавливает окружение сервиса
func WithEnvironment(environment string) Option {
	return func(c *Config) {
		c.environment = environment
	}
}

// WithVersion устанавливает версию сервиса
func WithVersion(version string) Option {
	return func(c *Config) {
		c.version = version
	}
}

// WithEnable включает/выключает метрики
func WithEnable(enable bool) Option {
	return func(c *Config) {
		c.enable = enable
	}
}

// WithNamespace устанавливает пространство имен для метрик
func WithNamespace(namespace string) Option {
	return func(c *Config) {
		c.namespace = namespace
	}
}

// WithAppName устанавливает имя приложения для метрик
func WithAppName(appName string) Option {
	return func(c *Config) {
		c.appName = appName
	}
}

// WithEndpoint устанавливает адрес OpenTelemetry коллектора
func WithEndpoint(endpoint string) Option {
	return func(c *Config) {
		c.endpoint = endpoint
	}
}

// WithTimeout устанавливает таймаут для операций с коллектором
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

// WithExportInterval устанавливает интервал экспорта метрик
func WithExportInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.exportInterval = interval
	}
}

// WithShutdownTimeout устанавливает таймаут для завершения работы
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.shutdownTimeout = timeout
	}
}
