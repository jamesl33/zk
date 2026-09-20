package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Open opens a sqlite3 database at the given path, configured for safe concurrent access from
// multiple 'zk' processes and goroutines: WAL mode lets readers and a writer proceed
// concurrently, and the busy timeout makes a writer wait for a lock instead of immediately
// failing with 'SQLITE_BUSY'. These are set via DSN parameters (rather than 'PRAGMA' after
// opening) so they apply to every connection database/sql opens from its pool, not just the
// first one.
func Open(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000", path)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}
