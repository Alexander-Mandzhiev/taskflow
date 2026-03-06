package tracing

import (
	"fmt"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var (
	_ contracts.TracingConfig = (*Config)(nil)
	_ contracts.Validatable   = (*Config)(nil)
)

type rawConfig struct {
	Enable   bool          `mapstructure:"enable"   env:"TRACING_ENABLE"`
	Endpoint string        `mapstructure:"endpoint" env:"TRACING_ENDPOINT"`
	Timeout  time.Duration `mapstructure:"timeout"  env:"TRACING_TIMEOUT"`

	// SampleRatio — процент сохраняемых трейсов (0–100). 100 = 100%, 10 = 10%. По умолчанию 10% для снижения нагрузки.
	SampleRatio float64 `mapstructure:"sample_ratio" env:"TRACING_SAMPLE_RATIO"`

	RetryEnabled         bool          `mapstructure:"retry_enabled"          env:"TRACING_RETRY_ENABLED"`
	RetryInitialInterval time.Duration `mapstructure:"retry_initial_interval" env:"TRACING_RETRY_INITIAL_INTERVAL"`
	RetryMaxInterval     time.Duration `mapstructure:"retry_max_interval"     env:"TRACING_RETRY_MAX_INTERVAL"`
	RetryMaxElapsedTime  time.Duration `mapstructure:"retry_max_elapsed_time" env:"TRACING_RETRY_MAX_ELAPSED_TIME"`

	EnableTraceContext bool          `mapstructure:"enable_trace_context" env:"TRACING_ENABLE_TRACE_CONTEXT"`
	EnableBaggage      bool          `mapstructure:"enable_baggage"        env:"TRACING_ENABLE_BAGGAGE"`
	ShutdownTimeout    time.Duration `mapstructure:"shutdown_timeout"     env:"TRACING_SHUTDOWN_TIMEOUT"`
}

// Config — конфиг модуля tracing (OpenTelemetry).
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		Enable:               false,
		Endpoint:             "localhost:4317",
		Timeout:              5 * time.Second,
		SampleRatio:          10, // 10%
		RetryEnabled:         true,
		RetryInitialInterval: 500 * time.Millisecond,
		RetryMaxInterval:     5 * time.Second,
		RetryMaxElapsedTime:  30 * time.Second,
		EnableTraceContext:   true,
		EnableBaggage:        true,
		ShutdownTimeout:      30 * time.Second,
	}
}

func (c *Config) Enable() bool                        { return c.raw.Enable }
func (c *Config) Endpoint() string                    { return c.raw.Endpoint }
func (c *Config) Timeout() time.Duration              { return c.raw.Timeout }
func (c *Config) SampleRatio() float64                { return c.raw.SampleRatio }
func (c *Config) RetryEnabled() bool                  { return c.raw.RetryEnabled }
func (c *Config) RetryInitialInterval() time.Duration { return c.raw.RetryInitialInterval }
func (c *Config) RetryMaxInterval() time.Duration     { return c.raw.RetryMaxInterval }
func (c *Config) RetryMaxElapsedTime() time.Duration  { return c.raw.RetryMaxElapsedTime }
func (c *Config) EnableTraceContext() bool            { return c.raw.EnableTraceContext }
func (c *Config) EnableBaggage() bool                 { return c.raw.EnableBaggage }
func (c *Config) ShutdownTimeout() time.Duration      { return c.raw.ShutdownTimeout }

// Validate проверяет корректность настроек трейсинга.
func (c *Config) Validate() error {
	if c.raw.SampleRatio < 0 || c.raw.SampleRatio > 100 {
		return fmt.Errorf("sample_ratio must be 0-100")
	}
	return nil
}
