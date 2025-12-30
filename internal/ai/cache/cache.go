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
type Cache[T any] struct {
	db    *sql.DB
	table string
}

// New - TODO
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
	  key integer unique,
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

// Get - TODO
func (c *Cache[T]) Get(ctx context.Context, prompt string) (*T, error) {
	hasher := crc32.NewIEEE()

	_, err := io.Copy(hasher, strings.NewReader(prompt))
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

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

	err = c.db.QueryRowContext(ctx, fmt.Sprintf(query, c.table), hasher.Sum32()).Scan(&result)

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
func (c *Cache[T]) Set(ctx context.Context, prompt string, result T) error {
	hasher := crc32.NewIEEE()

	_, err := io.Copy(hasher, strings.NewReader(prompt))
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// insert the embedding into the index
	const insert = `
	INSERT OR REPLACE INTO
	  %s
	VALUES
	  (?, ?);
	`

	_, err = c.db.ExecContext(
		ctx,
		fmt.Sprintf(insert, c.table),
		hasher.Sum32(),
		result,
	)
	if err != nil {
		return fmt.Errorf("failed to insert row: %w", err)
	}

	return nil
}
