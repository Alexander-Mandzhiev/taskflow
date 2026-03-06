package helpers

import (
	"os"
	"testing"
)

func TestResolveConfigPath(t *testing.T) {
	const defaultPath = "config/local.yaml"

	// Сохраняем и восстанавливаем CONFIG_PATH, чтобы не ломать другие тесты
	orig := os.Getenv("CONFIG_PATH")
	t.Cleanup(func() {
		var err error
		if orig == "" {
			err = os.Unsetenv("CONFIG_PATH")
		} else {
			err = os.Setenv("CONFIG_PATH", orig)
		}
		if err != nil {
			t.Fatalf("failed to restore CONFIG_PATH: %v", err)
		}
	})

	t.Run("default when CONFIG_PATH empty", func(t *testing.T) {
		if err := os.Unsetenv("CONFIG_PATH"); err != nil {
			t.Fatalf("os.Unsetenv: %v", err)
		}
		if got := ResolveConfigPath(defaultPath); got != defaultPath {
			t.Errorf("ResolveConfigPath(%q) = %q, want %q", defaultPath, got, defaultPath)
		}
	})

	t.Run("CONFIG_PATH overrides default", func(t *testing.T) {
		if err := os.Setenv("CONFIG_PATH", "/etc/app/production.yaml"); err != nil {
			t.Fatalf("os.Setenv: %v", err)
		}
		if got := ResolveConfigPath(defaultPath); got != "/etc/app/production.yaml" {
			t.Errorf("ResolveConfigPath(%q) = %q, want /etc/app/production.yaml", defaultPath, got)
		}
	})

	t.Run("CONFIG_PATH trimmed", func(t *testing.T) {
		if err := os.Setenv("CONFIG_PATH", "  /path/with/spaces.yaml  "); err != nil {
			t.Fatalf("os.Setenv: %v", err)
		}
		got := ResolveConfigPath(defaultPath)
		want := "/path/with/spaces.yaml"
		if got != want {
			t.Errorf("ResolveConfigPath(...) = %q, want %q", got, want)
		}
	})
}
