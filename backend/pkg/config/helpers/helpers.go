// Package helpers содержит общие функции для загрузки конфигурации:
// инициализация Viper из YAML (с подстановкой ${VAR} через os.ExpandEnv),
// получение секций, разрешение пути к файлу (ENV CONFIG_PATH или default).
//
// InitViper обычно вызывается один раз при старте приложения. Повторный вызов
// (например, для перезагрузки конфига по SIGHUP) возможен, но читатели GetSection
// в этот момент могут получить старое или новое состояние. Для горячей перезагрузки
// предпочтительна атомарная замена указателя (например, atomic.Pointer[*viper.Viper]).
package helpers

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type configState struct {
	viper *viper.Viper
	mu    sync.RWMutex
}

var globalState = &configState{}

// InitViper читает YAML по path, подставляет переменные окружения (os.ExpandEnv)
// и инициализирует Viper. При path == "" переводит в ENV-only режим (viper = nil).
func InitViper(path string) error {
	globalState.mu.Lock()
	defer globalState.mu.Unlock()

	if path == "" {
		globalState.viper = nil
		return nil
	}

	cleanPath := filepath.Clean(path)
	if !strings.HasSuffix(cleanPath, ".yaml") && !strings.HasSuffix(cleanPath, ".yml") {
		return fmt.Errorf("only YAML config files are allowed: %s", cleanPath)
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		globalState.viper = nil
		return fmt.Errorf("read config file %s: %w", cleanPath, err)
	}

	expanded := os.ExpandEnv(string(data))
	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewReader([]byte(expanded))); err != nil {
		globalState.viper = nil
		return fmt.Errorf("parse config file %s: %w", cleanPath, err)
	}

	globalState.viper = v
	return nil
}

// GetSection возвращает подсекцию по имени (например "mysql", "redis"). nil если viper не инициализирован или секции нет.
func GetSection(sectionName string) *viper.Viper {
	globalState.mu.RLock()
	defer globalState.mu.RUnlock()
	if globalState.viper == nil {
		return nil
	}
	return globalState.viper.Sub(sectionName)
}

// Reset сбрасывает состояние (только для тестов).
func Reset() {
	globalState.mu.Lock()
	defer globalState.mu.Unlock()
	globalState.viper = nil
}
