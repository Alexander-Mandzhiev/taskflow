package metric

import (
	"fmt"
	"time"

	"mkk/pkg/config/contracts"
)

var _ contracts.MetricConfig = (*Config)(nil)
var _ contracts.Validatable = (*Config)(nil)

type rawConfig struct {
	Enable           bool          `mapstructure:"enable" env:"METRIC_ENABLE"`
	Endpoint         string        `mapstructure:"endpoint" env:"METRIC_ENDPOINT"`
	Timeout          time.Duration `mapstructure:"timeout" env:"METRIC_TIMEOUT"`
	Namespace        string        `mapstructure:"namespace" env:"METRIC_NAMESPACE"`
	AppName          string        `mapstructure:"app_name" env:"METRIC_APP_NAME"`
	ExportInterval   time.Duration `mapstructure:"export_interval" env:"METRIC_EXPORT_INTERVAL"`
	ShutdownTimeout  time.Duration `mapstructure:"shutdown_timeout" env:"METRIC_SHUTDOWN_TIMEOUT"`
	BucketBoundaries []float64     `mapstructure:"bucket_boundaries" env:"METRIC_BUCKET_BOUNDARIES" envSeparator:","`
}

// Config — конфиг модуля metric (OpenTelemetry).
type Config struct {
	raw rawConfig
}

func defaultConfig() rawConfig {
	return rawConfig{
		Enable:          false,
		Endpoint:        "localhost:4317",
		Timeout:         30 * time.Second,
		Namespace:       "mkk",
		AppName:         "mkk",
		ExportInterval:  5 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		BucketBoundaries: []float64{
			0.0001, 0.0002, 0.0004, 0.0008, 0.0016, 0.0032, 0.0064, 0.0128,
			0.0256, 0.0512, 0.1024, 0.2048, 0.4096, 0.8192, 1.6384, 3.2768,
		},
	}
}

func (c *Config) Enable() bool                   { return c.raw.Enable }
func (c *Config) Endpoint() string               { return c.raw.Endpoint }
func (c *Config) Timeout() time.Duration         { return c.raw.Timeout }
func (c *Config) Namespace() string              { return c.raw.Namespace }
func (c *Config) AppName() string                { return c.raw.AppName }
func (c *Config) ExportInterval() time.Duration  { return c.raw.ExportInterval }
func (c *Config) ShutdownTimeout() time.Duration { return c.raw.ShutdownTimeout }
func (c *Config) BucketBoundaries() []float64 { return c.raw.BucketBoundaries }

// Validate проверяет корректность настроек метрик.
func (c *Config) Validate() error {
	if c.raw.ExportInterval <= 0 {
		return fmt.Errorf("export_interval must be positive")
	}
	return nil
}
