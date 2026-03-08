package testcontainers

import (
	"os"
	"strconv"
)

// TestConfigPath — путь к тестовому YAML относительно рабочей директории backend/.
const TestConfigPath = "config/test.yaml"

// TestEnv — параметры подключения к тестовым контейнерам (host/port от testcontainers).
// User, password, database и redis password берутся из config/test.yaml при config.Load().
// BackendURL — base URL API для apiclient (например http://localhost:4000), заполняется Setup().
type TestEnv struct {
	MySQLHost  string
	MySQLPort  int
	DSN        string
	RedisAddr  string
	BackendURL string // при полном стеке — baseURL для apiclient
}

// ApplyTestEnv выставляет в os.Environ только адреса контейнеров и CONFIG_PATH.
// Остальные переменные (MYSQL_USER, MYSQL_PASSWORD, MYSQL_DATABASE, REDIS_PASSWORD)
// pkg/config подхватит из config/test.yaml при Load().
func ApplyTestEnv(env *TestEnv) {
	_ = os.Setenv("CONFIG_PATH", TestConfigPath)
	_ = os.Setenv("MYSQL_HOST", env.MySQLHost)
	_ = os.Setenv("MYSQL_PORT", strconv.Itoa(env.MySQLPort))
	_ = os.Setenv("REDIS_ADDR", env.RedisAddr)
}

// ClearTestEnv снимает переменные тестового окружения (опционально, для изоляции тестов).
func ClearTestEnv() {
	_ = os.Unsetenv("CONFIG_PATH")
	_ = os.Unsetenv("MYSQL_HOST")
	_ = os.Unsetenv("MYSQL_PORT")
	_ = os.Unsetenv("REDIS_ADDR")
}
