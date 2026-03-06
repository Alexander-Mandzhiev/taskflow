package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var (
	_ contracts.RedisConfig = (*Config)(nil)
	_ contracts.Validatable = (*Config)(nil)
)

type rawConfig struct {
	Addr     string  `mapstructure:"addr"     env:"REDIS_ADDR"`
	Password string  `mapstructure:"password" env:"REDIS_PASSWORD"` //nolint:gosec // конфиг хранит только ссылку на секрет (env), не логируем и не сериализуем наружу
	Pool     rawPool `mapstructure:"pool"     envPrefix:"REDIS_POOL_"`
}

// Config — конфиг модуля redis (одна нода + пул, без кластера).
type Config struct {
	raw        rawConfig
	poolConfig *Pool
}

func defaultConfig() rawConfig {
	return rawConfig{
		Addr:     "localhost:6379",
		Password: "",
		Pool:     defaultPool(),
	}
}

// Addr возвращает адрес Redis (host:port). Если порт не указан, подставляется :6379.
func (c *Config) Addr() string {
	addr := strings.TrimSpace(c.raw.Addr)
	if addr == "" {
		return "localhost:6379"
	}
	if !strings.Contains(addr, ":") {
		return addr + ":6379"
	}
	return addr
}

func (c *Config) Password() string { return c.raw.Password }

// Timeout возвращает таймаут подключения из пула (удобно для cache.NewClient и т.п.).
func (c *Config) Timeout() time.Duration { return c.Pool().ConnTimeout() }

func (c *Config) Pool() contracts.RedisPoolConfig {
	return c.poolConfig
}

// Validate проверяет корректность настроек Redis (сырые поля, без геттеров).
func (c *Config) Validate() error {
	p := &c.raw.Pool
	if p.MaxActive < 0 || p.MaxIdle < 0 {
		return fmt.Errorf("pool max_active and max_idle must be >= 0")
	}
	if p.MaxActive > 0 && p.MaxIdle > p.MaxActive {
		return fmt.Errorf("pool max_idle (%d) cannot exceed max_active (%d)", p.MaxIdle, p.MaxActive)
	}
	if p.ConnTimeout < 0 || p.ReadTimeout < 0 || p.WriteTimeout < 0 || p.PoolTimeout < 0 || p.IdleTimeout < 0 {
		return fmt.Errorf("pool timeouts must be >= 0")
	}
	return nil
}
