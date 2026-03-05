package config

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/helpers"
	appmodule "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/internal/app"
	corsmodule "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/internal/cors"
	httpmodule "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/internal/http"
	loggermodule "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/internal/logger"
	metricmodule "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/internal/metric"
	mysqlmodule "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/internal/mysql"
	redismodule "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/internal/redis"
	sessionmodule "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/internal/session"
	tracingmodule "github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/internal/tracing"
)

// config реализует contracts.Provider через модульную сборку.
type config struct {
	appConfig     contracts.AppConfig
	httpConfig    contracts.HTTPConfig
	corsConfig    contracts.CORSConfig
	mysqlConfig   contracts.MySQLConfig
	redisConfig   contracts.RedisConfig
	sessionConfig contracts.SessionConfig
	loggerConfig  contracts.LoggerConfig
	tracingConfig contracts.TracingConfig
	metricConfig  contracts.MetricConfig
}

// buildModularConfig создаёт конфигурацию: каждый модуль загружается через свой New() (Defaults → YAML → ENV).
func buildModularConfig() (*config, error) {
	appCfg, err := appmodule.New()
	if err != nil {
		return nil, fmt.Errorf("app: %w", err)
	}
	httpCfg, err := httpmodule.New()
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	corsCfg, err := corsmodule.New()
	if err != nil {
		return nil, fmt.Errorf("cors: %w", err)
	}
	mysqlCfg, err := mysqlmodule.New()
	if err != nil {
		return nil, fmt.Errorf("mysql: %w", err)
	}
	redisCfg, err := redismodule.New()
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}
	sessionCfg, err := sessionmodule.New()
	if err != nil {
		return nil, fmt.Errorf("session: %w", err)
	}
	loggerCfg, err := loggermodule.New()
	if err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}
	tracingCfg, err := tracingmodule.New()
	if err != nil {
		return nil, fmt.Errorf("tracing: %w", err)
	}
	metricCfg, err := metricmodule.New()
	if err != nil {
		return nil, fmt.Errorf("metric: %w", err)
	}

	// Валидация модулей, реализующих Validatable. При добавлении нового модуля — добавить в этот слайс.
	modules := []struct {
		name string
		cfg  interface{}
	}{
		{"app", appCfg}, {"http", httpCfg}, {"cors", corsCfg}, {"mysql", mysqlCfg}, {"redis", redisCfg},
		{"session", sessionCfg}, {"logger", loggerCfg}, {"tracing", tracingCfg}, {"metric", metricCfg},
	}
	for _, m := range modules {
		if v, ok := m.cfg.(contracts.Validatable); ok {
			if err := v.Validate(); err != nil {
				return nil, fmt.Errorf("%s: %w", m.name, err)
			}
		}
	}
	return &config{
		appConfig:     appCfg,
		httpConfig:    httpCfg,
		corsConfig:    corsCfg,
		mysqlConfig:   mysqlCfg,
		redisConfig:   redisCfg,
		sessionConfig: sessionCfg,
		loggerConfig:  loggerCfg,
		tracingConfig: tracingCfg,
		metricConfig:  metricCfg,
	}, nil
}

// Load инициализирует Viper из YAML и собирает модульную конфигурацию.
// Путь к файлу задаётся в helpers.ResolveConfigPath: --config > CONFIG_PATH > ./config/development.yaml.
// Парсинг флагов выполняется при первом вызове (внутри helpers).
// Если файл не найден — ENV-only режим. Если файл найден, но невалиден — ошибка.
func Load(ctx context.Context) (contracts.Provider, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	path := helpers.ResolveConfigPath("./config/development.yaml")
	if err := helpers.InitViper(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("load config: %w", err)
		}
	}
	cfg, err := buildModularConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *config) App() contracts.AppConfig    { return c.appConfig }
func (c *config) HTTP() contracts.HTTPConfig   { return c.httpConfig }
func (c *config) CORS() contracts.CORSConfig   { return c.corsConfig }
func (c *config) MySQL() contracts.MySQLConfig { return c.mysqlConfig }
func (c *config) Redis() contracts.RedisConfig    { return c.redisConfig }
func (c *config) Session() contracts.SessionConfig { return c.sessionConfig }
func (c *config) Logger() contracts.LoggerConfig   { return c.loggerConfig }
func (c *config) Tracing() contracts.TracingConfig { return c.tracingConfig }
func (c *config) Metric() contracts.MetricConfig    { return c.metricConfig }
