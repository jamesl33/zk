package cache

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
)

// Cache defines a generic cache which is backed by a sqlite3 database.
//
// TODO (jamesl33): Add expiration.
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
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// create the table if it doesn't already exist.
	const create = `
	CREATE table IF NOT EXISTS %s (
	  key blob unique,
	  value blob
	);
	`

	_, err = db.ExecContext(ctx, fmt.Sprintf(create, table))
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	cache := Cache[T]{
		db:    db,
		table: table,
	}

	return &cache, nil
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

	err := c.db.QueryRowContext(ctx, fmt.Sprintf(query, c.table), key).Scan(&result)

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
		fmt.Sprintf(insert, c.table),
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
