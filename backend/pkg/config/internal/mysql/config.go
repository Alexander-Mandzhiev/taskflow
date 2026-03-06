package mysql

import (
	"fmt"
	"net"
	"strconv"
	"time"

	mysqldriver "github.com/go-sql-driver/mysql"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/contracts"
)

var (
	_ contracts.MySQLConfig = (*Config)(nil)
	_ contracts.Validatable = (*Config)(nil)
)

// rawConnection — параметры подключения. TODO: Add TLS support if moving to Public Cloud.
type rawConnection struct {
	Host     string `mapstructure:"host"     env:"HOST"`
	Port     int    `mapstructure:"port"     env:"PORT"`
	User     string `mapstructure:"user"     env:"USER"`
	Password string `mapstructure:"password" env:"PASSWORD"` //nolint:gosec // конфиг хранит только ссылку на секрет (env), не сам секрет; поле не сериализуется наружу
	Database string `mapstructure:"database" env:"DATABASE"`
}

type rawPool struct {
	MaxOpenConns    int           `mapstructure:"max_open_conns"    env:"MAX_OPEN_CONNS"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"    env:"MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" env:"CONN_MAX_LIFETIME"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" env:"CONN_MAX_IDLE_TIME"`
}

type rawConfig struct {
	Connection rawConnection `mapstructure:"connection" envPrefix:"MYSQL_"`
	Pool       rawPool       `mapstructure:"pool" envPrefix:"MYSQL_"`
}

// Config — конфиг модуля mysql (connection + pool).
type Config struct {
	raw rawConfig
	dsn string
}

func defaultConfig() rawConfig {
	return rawConfig{
		Connection: rawConnection{
			Host:     "localhost",
			Port:     3306,
			User:     "root",
			Password: "",
			Database: "mkk",
		},
		Pool: rawPool{
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
			ConnMaxIdleTime: 3 * time.Minute,
		},
	}
}

// DSN возвращает строку подключения MySQL (экранирование учётных данных через драйвер).
func (c *Config) DSN() string {
	return c.dsn
}

func (c *Config) MaxOpenConns() int              { return c.raw.Pool.MaxOpenConns }
func (c *Config) MaxIdleConns() int              { return c.raw.Pool.MaxIdleConns }
func (c *Config) ConnMaxLifetime() time.Duration { return c.raw.Pool.ConnMaxLifetime }
func (c *Config) ConnMaxIdleTime() time.Duration { return c.raw.Pool.ConnMaxIdleTime }

// Validate проверяет корректность настроек MySQL.
func (c *Config) Validate() error {
	if c.raw.Pool.MaxOpenConns < 0 {
		return fmt.Errorf("max_open_conns must be >= 0")
	}
	if c.raw.Pool.MaxIdleConns < 0 {
		return fmt.Errorf("max_idle_conns must be >= 0")
	}
	if c.raw.Pool.MaxOpenConns > 0 && c.raw.Pool.MaxIdleConns > c.raw.Pool.MaxOpenConns {
		return fmt.Errorf("max_idle_conns (%d) cannot be greater than max_open_conns (%d)",
			c.raw.Pool.MaxIdleConns, c.raw.Pool.MaxOpenConns)
	}
	if c.raw.Pool.ConnMaxLifetime < 0 {
		return fmt.Errorf("conn_max_lifetime must be >= 0")
	}
	if c.raw.Pool.ConnMaxIdleTime < 0 {
		return fmt.Errorf("conn_max_idle_time must be >= 0")
	}
	if c.raw.Connection.Port < 1 || c.raw.Connection.Port > 65535 {
		return fmt.Errorf("port must be 1-65535")
	}
	return nil
}

// buildDSN формирует DSN через mysql.Config (экранирует user/password).
func (c *Config) buildDSN() string {
	addr := net.JoinHostPort(c.raw.Connection.Host, strconv.Itoa(c.raw.Connection.Port))
	cfg := &mysqldriver.Config{
		User:      c.raw.Connection.User,
		Passwd:    c.raw.Connection.Password,
		Net:       "tcp",
		Addr:      addr,
		DBName:    c.raw.Connection.Database,
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       time.Local,
	}
	return cfg.FormatDSN()
}
