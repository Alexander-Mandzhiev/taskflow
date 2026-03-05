package config

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"

	"github.com/Alexander-Mandzhiev/taskflow/backend/pkg/config/helpers"
)

// Тесты загрузки конфига используют config/test.yaml (путь задаётся через CONFIG_PATH в каждом тесте).
// Не используйте t.Parallel() для тестов, вызывающих Load() или helpers — возможны гонки.

// testConfigPath возвращает путь к config/test.yaml.
// При запуске go test рабочая директория — каталог пакета (pkg/config), поэтому путь вверх к backend/config/.
func testConfigPath(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..", "config", "test.yaml")
}

func TestLoad_FromTestYAML(t *testing.T) {
	helpers.Reset()
	t.Cleanup(func() { helpers.Reset() })

	t.Setenv("CONFIG_PATH", testConfigPath(t))
	provider, err := Load(context.Background())
	if err != nil {
		t.Fatalf("Load from test.yaml: %v", err)
	}

	if name := provider.App().Name(); name != "testapp" {
		t.Errorf("App().Name() = %q, want testapp", name)
	}
	if env := provider.App().Environment(); env != "test" {
		t.Errorf("App().Environment() = %q, want test", env)
	}

	dsn := provider.MySQL().DSN()
	if dsn == "" {
		t.Error("MySQL().DSN() is empty")
	}
	if !strings.Contains(dsn, "testhost") || !strings.Contains(dsn, "3307") || !strings.Contains(dsn, "testdb") {
		t.Errorf("MySQL().DSN() = %q, expected host/port/db from test.yaml", dsn)
	}

	if addr := provider.Redis().Addr(); addr != "testredis:6380" {
		t.Errorf("Redis().Addr() = %q, want testredis:6380", addr)
	}

	if level := provider.Logger().Level(); level != "debug" {
		t.Errorf("Logger().Level() = %q, want debug", level)
	}

	if ratio := provider.Tracing().SampleRatio(); ratio != 50.5 {
		t.Errorf("Tracing().SampleRatio() = %v, want 50.5", ratio)
	}
}

func TestLoad_ENVOverridesYAML(t *testing.T) {
	helpers.Reset()
	t.Cleanup(func() { helpers.Reset() })

	t.Setenv("CONFIG_PATH", testConfigPath(t))
	t.Setenv("APP_NAME", "overridden")
	t.Setenv("LOGGER_LEVEL", "warn")

	provider, err := Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if name := provider.App().Name(); name != "overridden" {
		t.Errorf("App().Name() = %q, want overridden (ENV overrides YAML)", name)
	}
	if level := provider.Logger().Level(); level != "warn" {
		t.Errorf("Logger().Level() = %q, want warn (ENV overrides YAML)", level)
	}
}

func TestLoad_RedisAddr_DefaultPort(t *testing.T) {
	helpers.Reset()
	t.Cleanup(func() { helpers.Reset() })

	t.Setenv("CONFIG_PATH", testConfigPath(t))
	t.Setenv("REDIS_ADDR", "redis-only-host")

	provider, err := Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if addr := provider.Redis().Addr(); addr != "redis-only-host:6379" {
		t.Errorf("Redis().Addr() = %q, want redis-only-host:6379", addr)
	}
}

func TestLoad_ValidationFails(t *testing.T) {
	helpers.Reset()
	t.Cleanup(func() { helpers.Reset() })

	t.Setenv("CONFIG_PATH", testConfigPath(t))
	t.Setenv("LOGGER_LEVEL", "invalid")

	_, err := Load(context.Background())
	if err == nil {
		t.Fatal("Load with invalid config expected to fail")
	}
	if !strings.Contains(err.Error(), "logger") {
		t.Errorf("error should mention logger: %v", err)
	}
}

func TestLoad_ValidationFails_InvalidSampleRatio(t *testing.T) {
	helpers.Reset()
	t.Cleanup(func() { helpers.Reset() })

	t.Setenv("CONFIG_PATH", testConfigPath(t))
	t.Setenv("TRACING_SAMPLE_RATIO", "150")

	_, err := Load(context.Background())
	if err == nil {
		t.Fatal("Load with sample_ratio > 100 expected to fail")
	}
	if !strings.Contains(err.Error(), "tracing") {
		t.Errorf("error should mention tracing: %v", err)
	}
}

func TestLoad_ValidationFails_MySQL_MaxIdleGreaterThanMaxOpen(t *testing.T) {
	helpers.Reset()
	t.Cleanup(func() { helpers.Reset() })

	t.Setenv("CONFIG_PATH", testConfigPath(t))
	t.Setenv("MYSQL_MAX_OPEN_CONNS", "5")
	t.Setenv("MYSQL_MAX_IDLE_CONNS", "10")

	_, err := Load(context.Background())
	if err == nil {
		t.Fatal("Load with max_idle_conns > max_open_conns expected to fail")
	}
	if !strings.Contains(err.Error(), "mysql") {
		t.Errorf("error should mention mysql: %v", err)
	}
}

func TestLoad_ValidationFails_Redis_MaxIdleGreaterThanMaxActive(t *testing.T) {
	helpers.Reset()
	t.Cleanup(func() { helpers.Reset() })

	t.Setenv("CONFIG_PATH", testConfigPath(t))
	t.Setenv("REDIS_POOL_MAX_ACTIVE", "5")
	t.Setenv("REDIS_POOL_MAX_IDLE", "10")

	_, err := Load(context.Background())
	if err == nil {
		t.Fatal("Load with max_idle > max_active expected to fail")
	}
	if !strings.Contains(err.Error(), "redis") {
		t.Errorf("error should mention redis: %v", err)
	}
}

// TestLoad_MySQL_DSN_Escaping проверяет, что пароль и пользователь со спецсимволами
// корректно экранируются в DSN (драйвер mysql.Config.FormatDSN).
func TestLoad_MySQL_DSN_Escaping(t *testing.T) {
	helpers.Reset()
	t.Cleanup(func() { helpers.Reset() })

	t.Setenv("CONFIG_PATH", testConfigPath(t))
	t.Setenv("MYSQL_USER", "u@ser")
	t.Setenv("MYSQL_PASSWORD", "p@ss#word:123")

	provider, err := Load(context.Background())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	dsn := provider.MySQL().DSN()
	if dsn == "" {
		t.Fatal("DSN is empty")
	}

	parsed, err := mysql.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("ParseDSN (драйвер не смог разобрать DSN): %v", err)
	}
	if parsed.User != "u@ser" {
		t.Errorf("parsed User = %q, want u@ser", parsed.User)
	}
	if parsed.Passwd != "p@ss#word:123" {
		t.Errorf("parsed Passwd = %q, want p@ss#word:123", parsed.Passwd)
	}
}
