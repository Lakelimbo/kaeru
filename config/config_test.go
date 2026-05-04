package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Lakelimbo/kaeru/config"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "config.yml")

	content := `
name: testing instance
`

	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write testing config file: %v", err)
	}

	cfg, err := config.Load(file)
	if err != nil {
		t.Fatalf("failed to load testing config file: %v", err)
	}

	if cfg.Name != "testing instance" {
		t.Fatalf("expected testing config name to be 'testing instance', got '%s'", cfg.Name)
	}
}

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	cfg, err := config.Load("dummy")
	if err != nil {
		t.Fatalf("failed to load dummy config file: %v", err)
	}

	if cfg.Name != "kaeru instance" {
		t.Fatalf("expected config name to be 'kaeru instance', got '%s'", cfg.Name)
	}
	if cfg.Environment != "development" {
		t.Fatalf("expected defaullt environment to be 'development', got '%s'", cfg.Environment)
	}

	if cfg.OutputDir != ".kaeru" {
		t.Fatalf("expected output dir name to be '.kaeru', got '%s'", cfg.OutputDir)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Fatalf("expected default server host to be 127.0.0.1, got %s", cfg.Server.Host)
	}
	if cfg.Server.Port != 4040 {
		t.Fatalf("expected default server port to be 4040, got %d", cfg.Server.Port)
	}

	if cfg.Database.Path != "./.kaeru/db/kaeru.db" {
		t.Fatalf("expected database path to be './.kaeru/db/kaeru.db', got '%s'", cfg.Database.Path)
	}
	if cfg.Database.WAL != true {
		t.Fatal("expected WAL mode to be enabled by default")
	}

	if cfg.Logging.Level != "debug" {
		t.Fatalf("expected logging level to be 'debug', got %s", cfg.Logging.Level)
	}
}

func TestDefaultConfigOverride(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "config.yml")

	content := `
name: testing instance
server:
  port: 9050
`

	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write testing overridden config file: %v", err)
	}

	cfg, err := config.Load(file)
	if err != nil {
		t.Fatalf("failed to load overridden config file: %v", err)
	}

	if cfg.Name != "testing instance" {
		t.Fatalf("expected name override to be 'testing instance', got '%s'", cfg.Name)
	}
	if cfg.Server.Port != 9050 {
		t.Fatalf("expected server port override to be 9050, got %d", cfg.Server.Port)
	}
	// should still return default host value
	if cfg.Server.Host != "127.0.0.1" {
		t.Fatalf("expected default server host to be 127.0.0.1, got %s", cfg.Server.Host)
	}
}
