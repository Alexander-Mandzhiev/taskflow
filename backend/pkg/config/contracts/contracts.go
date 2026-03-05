package contracts

import (
	"time"
)

// Validatable — конфиг, который можно провалидировать при загрузке.
type Validatable interface {
	Validate() error
}

// Provider — единая точка доступа к конфигурации приложения.
// Каждый метод возвращает конфиг соответствующего модуля.
type Provider interface {
	App() AppConfig
	HTTP() HTTPConfig
	CORS() CORSConfig
	MySQL() MySQLConfig
	Redis() RedisConfig
	Session() SessionConfig
	Logger() LoggerConfig
	Tracing() TracingConfig
	Metric() MetricConfig
}

// AppConfig — базовые настройки приложения (имя, окружение, версия).
type AppConfig interface {
	Name() string
	Environment() string
	Version() string
}

// CORSConfig — настройки CORS для HTTP-сервера.
type CORSConfig interface {
	AllowedOrigins() []string
	AllowedMethods() []string
	AllowedHeaders() []string
	ExposedHeaders() []string
	AllowCredentials() bool
	MaxAge() int
}

// HTTPConfig — настройки HTTP-сервера (адрес, таймауты запроса/сервера).
type HTTPConfig interface {
	Address() string
	Timeout() time.Duration
	ReadHeaderTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	IdleTimeout() time.Duration
	MaxHeaderBytes() int
	ShutdownTimeout() time.Duration
}

// MySQLConfig — настройки подключения к MySQL (параметры подключения и пула).
type MySQLConfig interface {
	DSN() string
	MaxOpenConns() int
	MaxIdleConns() int
	ConnMaxLifetime() time.Duration
	ConnMaxIdleTime() time.Duration
}

// RedisConfig — настройки Redis (одна нода + пул соединений, без кластера).
type RedisConfig interface {
	Addr() string
	Password() string
	Timeout() time.Duration // удобный доступ к таймауту операций (равен Pool().ConnTimeout())
	Pool() RedisPoolConfig
}

// RedisPoolConfig — настройки пула соединений Redis.
type RedisPoolConfig interface {
	ConnTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	PoolTimeout() time.Duration
	MaxActive() int
	MaxIdle() int
	IdleTimeout() time.Duration
}

// SessionConfig — настройки сессий и cookie.
type SessionConfig interface {
	TTL() time.Duration
	IsSecure() bool
	CookieDomain() string
}

// LoggerConfig — настройки логгера (уровень, OTLP).
type LoggerConfig interface {
	Level() string
	AsJSON() bool
	Name() string
	Environment() string
	OTLPEnable() bool
	OTLPEndpoint() string
	OTLPShutdownTimeout() time.Duration
}

// TracingConfig — настройки OpenTelemetry трейсинга.
type TracingConfig interface {
	Enable() bool
	Endpoint() string
	Timeout() time.Duration
	SampleRatio() float64
	RetryEnabled() bool
	RetryInitialInterval() time.Duration
	RetryMaxInterval() time.Duration
	RetryMaxElapsedTime() time.Duration
	EnableTraceContext() bool
	EnableBaggage() bool
	ShutdownTimeout() time.Duration
}

// MetricConfig — настройки OpenTelemetry метрик.
type MetricConfig interface {
	Enable() bool
	Endpoint() string
	Timeout() time.Duration
	Namespace() string
	AppName() string
	ExportInterval() time.Duration
	ShutdownTimeout() time.Duration
	BucketBoundaries() []float64
}
