package database

import (
	"database/sql"

	"github.com/Lakelimbo/kaeru/internal/database/migrations"
	"github.com/Lakelimbo/kaeru/internal/logger"
	"github.com/pressly/goose/v3"
)

// migrate will take the SQL migrations (in /migrations) and
// use Goose to apply them.
//
// (TO-DO): have a down migrations function, or just ignore down
// migrations completely?
func migrate(db *sql.DB) error {
	if err := goose.SetDialect("sqlite3"); err != nil {
		logger.Fatalf("Failed to set dialect for migrations: %v", err)
	}

	goose.SetBaseFS(migrations.Embed)
	// goose.Up() needs to have the directory set to ".", because
	// it is understanding the context from migrations.Embed, not
	// the current directory of this file specifically.
	if err := goose.Up(db, "."); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}

	return nil
}
