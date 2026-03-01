package logger

import (
	"time"
)

// Option функциональная опция для конфигурации логгера
type Option func(*Config)

// Config содержит все настройки логгера
type Config struct {
	level       string
	asJSON      bool
	name        string
	environment string
	otlp        OTLPConfig
}

// OTLPConfig настройки для OpenTelemetry логов
type OTLPConfig struct {
	enable          bool
	endpoint        string
	shutdownTimeout time.Duration
}

// defaultConfig возвращает конфигурацию по умолчанию
func defaultConfig() *Config {
	return &Config{
		level:       "info",
		asJSON:      true,
		name:        "service",
		environment: "development",
		otlp: OTLPConfig{
			enable:          false,
			endpoint:        "localhost:4317",
			shutdownTimeout: 30 * time.Second,
		},
	}
}

// WithLevel устанавливает уровень логирования
func WithLevel(level string) Option {
	return func(c *Config) {
		c.level = level
	}
}

// WithJSON включает/выключает JSON формат
func WithJSON(asJSON bool) Option {
	return func(c *Config) {
		c.asJSON = asJSON
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

// WithOTLPEnable включает/выключает отправку логов в OpenTelemetry коллектор
func WithOTLPEnable(enable bool) Option {
	return func(c *Config) {
		c.otlp.enable = enable
	}
}

// WithOTLPEndpoint устанавливает адрес OpenTelemetry коллектора
func WithOTLPEndpoint(endpoint string) Option {
	return func(c *Config) {
		c.otlp.endpoint = endpoint
	}
}

// WithOTLPTimeout устанавливает таймаут для shutdown OTLP
func WithOTLPTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.otlp.shutdownTimeout = timeout
	}
}
