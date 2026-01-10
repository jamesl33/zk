package vector

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"iter"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jamesl33/zk/internal/ai"
	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
)

// TODO (jamesl33): Make this safe across instances of 'zk'.

// Enable SQLite vector search
func init() {
	sqlite_vec.Auto()
}

// DB exposes an API to index/find notes using SQLite vector search.
type DB struct {
	client *ai.Client
	db     *sql.DB
}

// New returns an initialized db.
func New(ctx context.Context, path string) (*DB, error) {
	client, err := ai.New(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// create the table if it doesn't already exist.
	const create = `
	CREATE table IF NOT EXISTS notes (
	  name text unique,
	  checksum integer,
	  embedding blob NOT NULL
	);
	`

	_, err = db.ExecContext(ctx, create)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	vector := DB{
		client: client,
		db:     db,
	}

	return &vector, nil
}

// Upsert adds the given note to the database (or updates it).
func (d *DB) Upsert(ctx context.Context, n *note.Note) error {
	checksum, err := n.Checksum()
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	skip, err := d.skip(ctx, n.Name(), checksum)
	if err != nil {
		return fmt.Errorf("failed to check if note requires updating: %w", err)
	}

	// Checksum matches, no further work required
	if skip {
		return nil
	}

	embedding, err := d.embed(ctx, n)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// No embedding returned, don't add to the index
	if embedding == nil {
		return nil
	}

	// insert the embedding into the index
	const insert = `
	INSERT OR REPLACE INTO
	  notes
	VALUES
	  (?, ?, ?);
	`

	_, err = d.db.ExecContext(
		ctx,
		insert,
		n.Name(),
		checksum,
		embedding,
	)
	if err != nil {
		return fmt.Errorf("failed to insert row: %w", err)
	}

	return nil
}

// Find some similar notes to the one provided.
func (d *DB) Find(ctx context.Context, n *note.Note) (iter.Seq2[*note.Note, error], error) {
	embedding, err := d.embed(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// No embedding, we can't find any similar notes
	if embedding == nil {
		return iterator.Empty[*note.Note](), nil
	}

	// query to find some similar notes
	const query = `
	SELECT
	  name,
      vec_distance_cosine (embedding, ?) as distance
	FROM
	  notes
	WHERE
	  name != ? AND distance <= 0.35
	ORDER BY
	  distance 
	LIMIT
	  5;
	`

	rows, err := d.db.QueryContext(ctx, query, embedding, n.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var (
		name     string
		distance float32
		names    = make([]string, 0)
	)

	for rows.Next() {
		err := rows.Scan(&name, &distance)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		names = append(names, name)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("unexpected error during iteration: %w", err)
	}

	// No notes, return an empty iterator.
	if len(names) == 0 {
		return iterator.Empty[*note.Note](), nil
	}

	matchers := hs.Map(names, func(n string) matcher.Matcher { return matcher.Name(n) })

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	return lister.Many(ctx), nil
}

// skip returns a boolean indicating whether we need to update the index entry.
func (d *DB) skip(ctx context.Context, name string, current uint32) (bool, error) {
	// query to acquire the existing checksum
	const query = `
	SELECT
	  checksum
	FROM
	  notes
	WHERE
	  name = ?
	`

	var indexed uint32

	err := d.db.QueryRowContext(ctx, query, name).Scan(&indexed)

	// Not found, we need to update
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to query for note checksum: %w", err)
	}

	return current == indexed, nil
}

// embed returns a vector embedding for the given note.
func (d *DB) embed(ctx context.Context, n *note.Note) ([]byte, error) {
	// Reduce false positives caused by embedding empty notes
	if len(n.Frontmatter.Tags) == 0 && len(n.Body) == 0 {
		return nil, nil
	}

	var input bytes.Buffer

	_, err := n.WriteTo(&input)
	if err != nil {
		return nil, fmt.Errorf("failed to write note to buffer: %w", err)
	}

	vec, err := d.client.Embed(ctx, input.String())
	if err != nil {
		return nil, fmt.Errorf("failed to embed buffer: %w", err)
	}

	// We didn't receive an embedding
	if len(vec) == 1 {
		return nil, nil
	}

	serial, err := sqlite_vec.SerializeFloat32(vec)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize embedding: %w", err)
	}

	return serial, nil
}

// Close frees resources used by the database.
func (d *DB) Close() error {
	return d.db.Close()
}
