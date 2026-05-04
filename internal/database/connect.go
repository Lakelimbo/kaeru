package database

import (
	"github.com/Lakelimbo/kaeru/config"
	"github.com/pocketbase/dbx"
	_ "modernc.org/sqlite"
)

// connect will open the SQLite database and also add pragmas
// (such as WAL) if enabled.
func connect(cfg *config.DatabaseConfig) (*dbx.DB, error) {
	pragmas := "?_pragma=busy_timeout(10000)"
	if cfg.WAL {
		pragmas += "&_pragma=journal_mode(WAL)"
	}
	pragmas += "&_pragma=journal_size_limit(200000000)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)&_pragma=temp_store(MEMORY)&_pragma=cache_size(-32000)"

	db, err := dbx.Open("sqlite", cfg.Path+pragmas)
	if err != nil {
		return nil, err
	}

	return db, nil
}
