package tests

import (
	"path/filepath"

	"github.com/Lakelimbo/kaeru/config"
	"github.com/Lakelimbo/kaeru/internal/containers"
	"github.com/Lakelimbo/kaeru/internal/database"
	"github.com/Lakelimbo/kaeru/internal/jobs"
	"github.com/Lakelimbo/kaeru/internal/server"
)

// NewTestApp creates a testing instance of Kaeru for
// integration tests.
func NewTestApp(cfg *config.Config) *server.App {
	db := database.New(cfg)
	docker, _ := containers.NewDockerRepository(cfg, db)
	messaging := jobs.NewPubSub()
	jobs := jobs.Spawn(db, docker, messaging)

	return server.New(cfg, db, docker, jobs)
}

// TempDB sets a basic path for the SQLite database
// for any integration tests that need it.
//
// dir will be likely a temporary dir (i.e created by
// t.TempDir() ).
func TempDB(dir string) *config.Config {
	return &config.Config{
		Database: config.DatabaseConfig{
			Path: filepath.Join(dir, "test.db"),
		},
	}
}
