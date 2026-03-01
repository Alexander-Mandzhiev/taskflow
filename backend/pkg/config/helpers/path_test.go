package helpers

import (
	"os"
	"testing"
)

func TestResolveConfigPath(t *testing.T) {
	const defaultPath = "config/development.yaml"

	// Сохраняем и восстанавливаем CONFIG_PATH, чтобы не ломать другие тесты
	orig := os.Getenv("CONFIG_PATH")
	t.Cleanup(func() {
		if orig == "" {
			os.Unsetenv("CONFIG_PATH")
		} else {
			os.Setenv("CONFIG_PATH", orig)
		}
	})

	t.Run("default when CONFIG_PATH empty", func(t *testing.T) {
		os.Unsetenv("CONFIG_PATH")
		if got := ResolveConfigPath(defaultPath); got != defaultPath {
			t.Errorf("ResolveConfigPath(%q) = %q, want %q", defaultPath, got, defaultPath)
		}
	})

	t.Run("CONFIG_PATH overrides default", func(t *testing.T) {
		os.Setenv("CONFIG_PATH", "/etc/app/production.yaml")
		if got := ResolveConfigPath(defaultPath); got != "/etc/app/production.yaml" {
			t.Errorf("ResolveConfigPath(%q) = %q, want /etc/app/production.yaml", defaultPath, got)
		}
	})

	t.Run("CONFIG_PATH trimmed", func(t *testing.T) {
		os.Setenv("CONFIG_PATH", "  /path/with/spaces.yaml  ")
		got := ResolveConfigPath(defaultPath)
		want := "/path/with/spaces.yaml"
		if got != want {
			t.Errorf("ResolveConfigPath(...) = %q, want %q", got, want)
		}
	})
}
