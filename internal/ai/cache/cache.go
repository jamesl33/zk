package cache

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jamesl33/zk/internal/sqlite"
)

// Cache defines a generic cache which is backed by a sqlite3 database.
type Cache[T any] struct {
	db    *sql.DB
	table string
}

// New creates a new cache.
func New[T any](
	ctx context.Context,
	path string,
	table string,
) (*Cache[T], error) {
	cache := Cache[T]{table: table}

	err := cache.setup(ctx, path)
	if err == nil {
		return &cache, nil
	}

	cache.Close() //nolint:errcheck

	return nil, err
}

// setup opens the database and creates the table. On error the cache may be partially set up, the caller must Close
// it.
func (c *Cache[T]) setup(ctx context.Context, path string) error {
	var err error

	c.db, err = sqlite.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// create the table if it doesn't already exist.
	const create = `
	CREATE table IF NOT EXISTS %s (
	  key blob unique,
	  value blob
	);
	`

	_, err = c.db.ExecContext(ctx, fmt.Sprintf(create, quote(c.table)))
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// Close frees resources used by the cache.
func (c *Cache[T]) Close() error {
	if c.db == nil {
		return nil
	}

	return c.db.Close()
}

// Get a value from the cache.
func (c *Cache[T]) Get(ctx context.Context, prompt string) (*T, error) {
	key := checksum(prompt)

	// query to acquire the existing prompt checksum
	const query = `
	SELECT
	  value
	FROM
	  %s
	WHERE
	  key = ?
	`

	var result T

	err := c.db.QueryRowContext(ctx, fmt.Sprintf(query, quote(c.table)), key).Scan(&result)

	// Not found, we need to update
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query row: %w", err)
	}

	return &result, nil
}

// Set a value in the cache.
func (c *Cache[T]) Set(ctx context.Context, prompt string, result T) error {
	key := checksum(prompt)

	// insert the embedding into the index
	const insert = `
	INSERT OR REPLACE INTO
	  %s
	VALUES
	  (?, ?);
	`

	_, err := c.db.ExecContext(
		ctx,
		fmt.Sprintf(insert, quote(c.table)),
		key,
		result,
	)
	if err != nil {
		return fmt.Errorf("failed to insert row: %w", err)
	}

	return nil
}

// checksum returns a fixed-size digest of the given prompt, used as the
// cache key. SHA-256 is used (rather than a shorter checksum like CRC32) so
// that a collision between two different prompts is computationally
// infeasible to find, rather than something a handful of cache entries
// could realistically hit by chance.
func checksum(prompt string) []byte {
	sum := sha256.Sum256([]byte(prompt))
	return sum[:]
}

// quote returns the given SQL identifier quoted, so arbitrary names (e.g. model names containing
// ':' or '-') are safe to interpolate.
func quote(identifier string) string {
	return `"` + strings.ReplaceAll(identifier, `"`, `""`) + `"`
}
