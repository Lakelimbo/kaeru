package database

import (
	"os"
	"path/filepath"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/pocketbase/dbx"
)

// The Database instance, which is controlled with dbx.DB.
type DB struct {
	*dbx.DB
}

func New(cfg *config.Config) *dbx.DB {
	dir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Fatalf("Failed to create database directory: %v", err)
	}

	db, err := connect(&cfg.Database)
	if err != nil {
		logger.Fatalf("Unable to open SQLite database: %v", err)
	}

	if err := db.DB().Ping(); err != nil {
		logger.Fatalf("Failed to ping database: %v", err)
	}

	if err := migrate(db.DB()); err != nil {
		return nil
	}

	db.LogFunc = logger.Debugf

	logger.Info("SQLite connection initialized")
	return db
}
