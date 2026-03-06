package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitViper(t *testing.T) {
	t.Cleanup(func() { Reset() })

	t.Run("empty path sets ENV-only mode", func(t *testing.T) {
		Reset()
		if err := InitViper(""); err != nil {
			t.Fatalf("InitViper(\"\"): %v", err)
		}
		if got := GetSection("app"); got != nil {
			t.Error("GetSection after InitViper(\"\") should return nil (ENV-only mode)")
		}
	})

	t.Run("rejects non-YAML extension", func(t *testing.T) {
		Reset()
		err := InitViper("/etc/app/config.json")
		if err == nil {
			t.Fatal("InitViper with .json expected to fail")
		}
		if !strings.Contains(err.Error(), "only YAML") {
			t.Errorf("error should mention YAML: %v", err)
		}
	})

	t.Run("returns error for missing file", func(t *testing.T) {
		Reset()
		err := InitViper(filepath.Join(t.TempDir(), "missing.yaml"))
		if err == nil {
			t.Fatal("InitViper with missing file expected to fail")
		}
		if !strings.Contains(err.Error(), "read config file") {
			t.Errorf("error should mention read: %v", err)
		}
	})

	t.Run("returns error for invalid YAML content", func(t *testing.T) {
		Reset()
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.yaml")
		if err := os.WriteFile(path, []byte("invalid: yaml: [[["), 0o600); err != nil {
			t.Fatalf("write test file: %v", err)
		}
		err := InitViper(path)
		if err == nil {
			t.Fatal("InitViper with invalid YAML expected to fail")
		}
		if !strings.Contains(err.Error(), "parse") {
			t.Errorf("error should mention parse: %v", err)
		}
	})

	t.Run("loads valid YAML and expands env", func(t *testing.T) {
		Reset()
		dir := t.TempDir()
		path := filepath.Join(dir, "app.yaml")
		// ${TEST_CONFIG_VAR} will be expanded
		yaml := "app:\n  name: ${TEST_CONFIG_VAR}\n  environment: test\n"
		if err := os.WriteFile(path, []byte(yaml), 0o600); err != nil {
			t.Fatalf("write test file: %v", err)
		}
		t.Setenv("TEST_CONFIG_VAR", "expanded-name")
		if err := InitViper(path); err != nil {
			t.Fatalf("InitViper: %v", err)
		}
		sub := GetSection("app")
		if sub == nil {
			t.Fatal("GetSection(\"app\") should not be nil after successful InitViper")
		}
		if got := sub.GetString("name"); got != "expanded-name" {
			t.Errorf("app.name = %q, want expanded-name (env expanded)", got)
		}
	})
}

func TestGetSection(t *testing.T) {
	t.Cleanup(func() { Reset() })

	t.Run("returns nil when viper not initialized", func(t *testing.T) {
		Reset()
		if got := GetSection("app"); got != nil {
			t.Errorf("GetSection before InitViper = %v, want nil", got)
		}
	})

	t.Run("returns sub when section exists", func(t *testing.T) {
		Reset()
		dir := t.TempDir()
		path := filepath.Join(dir, "cfg.yaml")
		yaml := "mysql:\n  host: localhost\n  port: 3306\n"
		if err := os.WriteFile(path, []byte(yaml), 0o600); err != nil {
			t.Fatalf("write test file: %v", err)
		}
		if err := InitViper(path); err != nil {
			t.Fatalf("InitViper: %v", err)
		}
		sub := GetSection("mysql")
		if sub == nil {
			t.Fatal("GetSection(\"mysql\") = nil")
		}
		if got := sub.GetString("host"); got != "localhost" {
			t.Errorf("mysql.host = %q, want localhost", got)
		}
		if got := sub.GetInt("port"); got != 3306 {
			t.Errorf("mysql.port = %v, want 3306", got)
		}
	})
}
