package database_test

import (
	"path/filepath"
	"testing"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/database"
)

func TestDatabaseInstance(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		Database: config.DatabaseConfig{
			Path: filepath.Join(t.TempDir(), "test.db"),
		},
	}

	db := database.New(&cfg)
	if db == nil {
		t.Fatal("expected database instance, got nil")
	}
}

func TestDBX(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		Database: config.DatabaseConfig{
			Path: filepath.Join(t.TempDir(), "test.db"),
		},
	}

	db := database.New(&cfg)
	if db == nil {
		t.Fatal("expected database instance, got nil")
	}
}
