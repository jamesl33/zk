package cache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"strings"
)

// Cache - TODO
//
// TODO (jamesl33): Add expiration.
type Cache struct {
	db *sql.DB
}

// New - TODO
func New(ctx context.Context, path string) (*Cache, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// create the table if it doesn't already exist.
	const create = `
	CREATE table IF NOT EXISTS cache (
	  prompt integer unique,
	  result text
	);
	`

	_, err = db.ExecContext(ctx, create)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	cache := Cache{
		db: db,
	}

	return &cache, nil
}

// Get - TODO
func (c *Cache) Get(ctx context.Context, prompt string) (*string, error) {
	hasher := crc32.NewIEEE()

	_, err := io.Copy(hasher, strings.NewReader(prompt))
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// query to acquire the existing prompt checksum
	const query = `
	SELECT
	  result
	FROM
	  cache
	WHERE
	  prompt = ?
	`

	var result string

	err = c.db.QueryRowContext(ctx, query, hasher.Sum32()).Scan(&result)

	// Not found, we need to update
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return &result, nil
}

// Set - TODO
func (c *Cache) Set(ctx context.Context, prompt, result string) error {
	hasher := crc32.NewIEEE()

	_, err := io.Copy(hasher, strings.NewReader(prompt))
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// insert the embedding into the index
	const insert = `
	INSERT OR REPLACE INTO
	  cache
	VALUES
	  (?, ?);
	`

	_, err = c.db.ExecContext(
		ctx,
		insert,
		hasher.Sum32(),
		result,
	)
	if err != nil {
		return fmt.Errorf("failed to insert row: %w", err)
	}

	return nil
}
