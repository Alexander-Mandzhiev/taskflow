package redis

import (
	"time"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var _ contracts.RedisPoolConfig = (*Pool)(nil)

type rawPool struct {
	MaxActive    int           `mapstructure:"max_active"    env:"MAX_ACTIVE"`
	MaxIdle      int           `mapstructure:"max_idle"     env:"MAX_IDLE"`
	ConnTimeout  time.Duration `mapstructure:"conn_timeout"  env:"CONN_TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" env:"WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" env:"IDLE_TIMEOUT"`
	PoolTimeout  time.Duration `mapstructure:"pool_timeout" env:"POOL_TIMEOUT"`
}

// Pool — настройки пула соединений Redis.
type Pool struct {
	raw rawPool
}

func defaultPool() rawPool {
	return rawPool{
		MaxActive:    10,
		MaxIdle:      5,
		ConnTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout:  240 * time.Second,
		PoolTimeout:  4 * time.Second,
	}
}

func (p *Pool) ConnTimeout() time.Duration  { return p.raw.ConnTimeout }
func (p *Pool) ReadTimeout() time.Duration  { return p.raw.ReadTimeout }
func (p *Pool) WriteTimeout() time.Duration { return p.raw.WriteTimeout }
func (p *Pool) PoolTimeout() time.Duration  { return p.raw.PoolTimeout }
func (p *Pool) MaxActive() int              { return p.raw.MaxActive }
func (p *Pool) MaxIdle() int                { return p.raw.MaxIdle }
func (p *Pool) IdleTimeout() time.Duration  { return p.raw.IdleTimeout }
