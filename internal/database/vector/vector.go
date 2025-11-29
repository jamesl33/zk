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

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes/lister"
	"github.com/jamesl33/zk/internal/notes/matcher"
	"github.com/ollama/ollama/api"
)

// DB - TODO
type DB struct {
	client *api.Client
	db     *sql.DB
}

// New - TODO
func New(ctx context.Context, path string) (*DB, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// create - TODO
	const create = `
	CREATE table IF NOT EXISTS notes (
	  name text unique,
	  checksum integer,
	  embedding blob NOT NULL
	);
	`

	_, err = db.ExecContext(ctx, create)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	vector := DB{
		client: client,
		db:     db,
	}

	return &vector, nil
}

// Upsert - TODO
func (d *DB) Upsert(ctx context.Context, n *note.Note) error {
	checksum, err := n.Checksum()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	skip, err := d.skip(ctx, n.Name(), checksum)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if skip {
		return nil
	}

	embedding, err := d.embed(ctx, n)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if embedding == nil {
		return nil
	}

	// insert - TODO
	const insert = `
	INSERT INTO
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
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// Find - TODO
func (d *DB) Find(ctx context.Context, n *note.Note) (iter.Seq2[*note.Note, error], error) {
	embedding, err := d.embed(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if embedding == nil {
		return iterator.Empty[*note.Note](), nil
	}

	// query - TODO
	const query = `
	SELECT
	  name,
      vec_distance_cosine (embedding, ?) as distance
	FROM
	  notes
	WHERE
	  name != ? AND distance <= 0.5
	ORDER BY
	  distance 
	LIMIT
	  5;
	`

	rows, err := d.db.QueryContext(ctx, query, embedding, n.Name())
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	var (
		name     string
		distance float32
		names    = make([]string, 0)
	)

	for rows.Next() {
		err := rows.Scan(&name, &distance)
		if err != nil {
			return nil, fmt.Errorf("%w", err) // TODO
		}

		names = append(names, name)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if len(names) == 0 {
		return iterator.Empty[*note.Note](), nil
	}

	matchers := hs.Map(names, func(n string) matcher.Matcher { return matcher.Name(n) })

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return lister.Many(ctx), nil
}

// skip - TODO
func (d *DB) skip(ctx context.Context, name string, current uint32) (bool, error) {
	// query - TODO
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

	// TODO
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("%w", err) // TODO
	}

	return current == indexed, nil
}

// embed - TODO
//
// TODO (jamesl33): Make the model configurable.
// TODO (jamesl33): Handle the 2k context window.
func (d *DB) embed(ctx context.Context, n *note.Note) ([]byte, error) {
	var input bytes.Buffer

	err := n.WriteTo(&input)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	req := api.EmbedRequest{
		Model: "embeddinggemma:300m",
		Input: input.String(),
	}

	resp, err := d.client.Embed(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if len(resp.Embeddings) != 1 {
		return nil, nil
	}

	serial, err := sqlite_vec.SerializeFloat32(resp.Embeddings[0])
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return serial, nil
}

// Close - TODO
func (d *DB) Close() error {
	return d.db.Close()
}
