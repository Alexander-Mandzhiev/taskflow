package tracing

import (
	"time"
)

// Option функциональная опция для конфигурации трейсинга
type Option func(*Config)

// Config содержит все настройки трейсинга
type Config struct {
	name                 string
	environment          string
	version              string
	enable               bool
	endpoint             string
	timeout              time.Duration
	sampleRatio          int
	retryEnabled         bool
	retryInitialInterval time.Duration
	retryMaxInterval     time.Duration
	retryMaxElapsedTime  time.Duration
	enableTraceContext   bool
	enableBaggage        bool
}

// defaultConfig возвращает конфигурацию по умолчанию
func defaultConfig() *Config {
	return &Config{
		name:                 "service",
		environment:          "development",
		version:              "1.0.0",
		enable:               false,
		endpoint:             "localhost:4317",
		timeout:              30 * time.Second,
		sampleRatio:          100,
		retryEnabled:         true,
		retryInitialInterval: 5 * time.Second,
		retryMaxInterval:     30 * time.Second,
		retryMaxElapsedTime:  2 * time.Minute,
		enableTraceContext:   true,
		enableBaggage:        true,
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

// WithEndpoint устанавливает адрес OpenTelemetry коллектора
func WithEndpoint(endpoint string) Option {
	return func(c *Config) {
		c.endpoint = endpoint
	}
}

// WithEnable включает/выключает трейсинг
func WithEnable(enable bool) Option {
	return func(c *Config) {
		c.enable = enable
	}
}

// WithTimeout устанавливает таймаут для операций с коллектором
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

// WithSampleRatio устанавливает процент сохраняемых трейсов (0 - 100)
func WithSampleRatio(ratio int) Option {
	return func(c *Config) {
		c.sampleRatio = ratio
	}
}

// WithRetryEnabled включает/выключает повторные попытки
func WithRetryEnabled(enabled bool) Option {
	return func(c *Config) {
		c.retryEnabled = enabled
	}
}

// WithRetryInitialInterval устанавливает начальный интервал между попытками
func WithRetryInitialInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.retryInitialInterval = interval
	}
}

// WithRetryMaxInterval устанавливает максимальный интервал между попытками
func WithRetryMaxInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.retryMaxInterval = interval
	}
}

// WithRetryMaxElapsedTime устанавливает максимальное время на все попытки
func WithRetryMaxElapsedTime(time time.Duration) Option {
	return func(c *Config) {
		c.retryMaxElapsedTime = time
	}
}

// WithTraceContext включает/выключает W3C TraceContext
func WithTraceContext(enable bool) Option {
	return func(c *Config) {
		c.enableTraceContext = enable
	}
}

// WithBaggage включает/выключает передачу дополнительного контекста
func WithBaggage(enable bool) Option {
	return func(c *Config) {
		c.enableBaggage = enable
	}
}
