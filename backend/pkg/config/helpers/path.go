package helpers

import (
	"flag"
	"os"
	"strings"
	"sync"
)

var (
	once       sync.Once
	configFlag string
)

// initFlags регистрирует флаг -config и парсит командную строку.
// Вызывается один раз при первом обращении к ResolveConfigPath.
// Парсинг флагов выполняется только здесь (main флаги не парсит до конфига).
func initFlags() {
	flag.StringVar(&configFlag, "config", "", "path to configuration file (YAML)")
	if !flag.Parsed() {
		flag.Parse()
	}
}

// ResolveConfigPath возвращает путь к конфигурационному файлу.
// Приоритет (от высшего к низшему):
//  1. Флаг командной строки -config
//  2. Переменная окружения CONFIG_PATH (пробелы обрезаются)
//  3. defaultPath
func ResolveConfigPath(defaultPath string) string {
	once.Do(initFlags)

	if configFlag != "" {
		return configFlag
	}
	if v := strings.TrimSpace(os.Getenv("CONFIG_PATH")); v != "" {
		return v
	}
	return defaultPath
}
