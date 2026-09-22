package vector

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jamesl33/zk/internal/ai"
	"github.com/jamesl33/zk/internal/chunker"
	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/sqlite"
	"github.com/jamesl33/zk/internal/vault"
)

// Enable SQLite vector search
func init() {
	sqlite_vec.Auto()
}

// DB exposes an API to index/find notes using SQLite vector search.
type DB struct {
	client ai.Client
	db     *sql.DB
}

// New returns an initialized db.
func New(ctx context.Context, path string) (*DB, error) {
	client, err := ai.NewOllama(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	db, err := sqlite.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	vector := DB{
		client: client,
		db:     db,
	}

	err = vector.init(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &vector, nil
}

// init the database by creating the required table.
func (d *DB) init(ctx context.Context) error {
	// create the table if it doesn't already exist. Notes may be split into multiple chunks (see chunker), so a note
	// can have more than one row, keyed by (name, chunk).
	const create = `
	CREATE table IF NOT EXISTS notes (
	  name text,
	  chunk integer,
	  checksum blob,
	  embedding blob NOT NULL,
	  PRIMARY KEY (name, chunk)
	);
	`

	_, err := d.db.ExecContext(ctx, create)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
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

	embeddings, err := d.embed(ctx, n)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	// Remove any existing rows for this note; the number of chunks may have changed since the last time it was indexed.
	_, err = tx.ExecContext(ctx, `DELETE FROM notes WHERE name = ?;`, n.Name())
	if err != nil {
		return fmt.Errorf("failed to delete stale rows: %w", err)
	}

	// No embeddings returned, don't add to the index
	if len(embeddings) == 0 {
		return tx.Commit()
	}

	err = d.upsert(ctx, tx, n.Name(), checksum, embeddings)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// upsert the embeddings.
func (d *DB) upsert(ctx context.Context, tx *sql.Tx, name string, checksum []byte, embeddings [][]byte) error {
	const insert = `
	INSERT OR REPLACE INTO
	  notes
	VALUES
	  (?, ?, ?, ?);
	`

	for i, embedding := range embeddings {
		_, err := tx.ExecContext(ctx, insert, name, i, checksum, embedding)
		if err != nil {
			return fmt.Errorf("failed to insert row for chunk %d: %w", i, err)
		}
	}

	return nil
}

// Find some similar notes to the one provided.
func (d *DB) Find(ctx context.Context, n *note.Note) ([]*note.Note, error) {
	embeddings, err := d.embed(ctx, n)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// No embedding, we can't find any similar notes
	if len(embeddings) == 0 {
		return make([]*note.Note, 0), nil
	}

	// If the note itself needed to be split into multiple chunks (rare -- notes are expected to stay small/atomic), use
	// the leading chunk as a single representative query vector rather than searching once per chunk and merging
	// results.
	embedding := embeddings[0]

	// query to find some similar notes; a note may have multiple chunk rows, so collapse to one distance per note using
	// its closest-matching chunk.
	const query = `
	SELECT
	  name,
      MIN(vec_distance_cosine (embedding, ?)) as distance
	FROM
	  notes
	WHERE
	  name != ?
	GROUP BY
	  name
	HAVING
	  distance <= 0.6
	ORDER BY
	  distance
	LIMIT 15
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
		return make([]*note.Note, 0), nil
	}

	root, err := vault.RootRel(".")
	if err != nil {
		return nil, fmt.Errorf("failed to find vault root: %w", err)
	}

	matchers := hs.Map(names, func(n string) matcher.Matcher { return matcher.Name(n) })

	lister, err := lister.NewLister(
		lister.WithPath(root),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	notes := make([]*note.Note, len(matchers))

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		notes[slices.IndexFunc(names, func(name string) bool { return name == n.Name() })] = n
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to collect notes: %w", err)
	}

	// We may have entries in the index, but not on disk (e.g. deleted notes)
	notes = slices.DeleteFunc(notes, func(n *note.Note) bool { return n == nil })

	// Remove unused capacity
	notes = slices.Clip(notes)

	return notes, nil
}

// skip returns a boolean indicating whether we need to update the index entry.
func (d *DB) skip(ctx context.Context, name string, current []byte) (bool, error) {
	// query to acquire the existing checksum; every chunk row for a name shares the same
	// checksum, so any one row will do.
	const query = `
	SELECT
	  checksum
	FROM
	  notes
	WHERE
	  name = ?
	LIMIT 1
	`

	var indexed []byte

	err := d.db.QueryRowContext(ctx, query, name).Scan(&indexed)

	// Not found, we need to update
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to query for note checksum: %w", err)
	}

	return bytes.Equal(current, indexed), nil
}

// embed returns one vector embedding per chunk for the given note. Notes that exceed the embedding model's context
// window are split into multiple chunks (see chunker); most notes produce exactly one.
func (d *DB) embed(ctx context.Context, n *note.Note) ([][]byte, error) {
	body, err := n.GetBody()
	if err != nil {
		return nil, fmt.Errorf("failed to get body: %w", err)
	}

	// Reduce false positives caused by embedding empty notes
	if len(n.Frontmatter.Tags) == 0 && len(body) == 0 {
		return nil, nil
	}

	var input bytes.Buffer

	_, err = n.WriteTo(&input)
	if err != nil {
		return nil, fmt.Errorf("failed to write note to buffer: %w", err)
	}

	// WriteTo always writes "---\n<yaml>---\n" then exactly one blank line then the body, so splitting on the first
	// blank line cleanly separates them without re-serializing Frontmatter.
	var (
		parts       = strings.SplitN(input.String(), "\n\n", 2)
		frontmatter = parts[0]
	)

	var bodyText string
	if len(parts) > 1 {
		bodyText = parts[1]
	}

	const (
		// context is the context window of the embedding model. Input beyond this errors, rather than being truncated.
		context = 2048

		// limit is the character budget for a single chunk; assumes ~3 chars/token for markdown, minus a 20% margin.
		limit = context * 3 * 4 / 5

		// mn is the smallest character budget attempted before giving up.
		mn = 256
	)

	// The character budget is only an estimate of the token count; dense content (e.g. identifiers, hashes or encoded
	// URLs) can exceed the context length regardless. Re-chunk the whole note with a halved budget until every chunk
	// fits.
	for limit := limit; ; limit /= 2 {
		embeddings, err := d.embedChunks(ctx, frontmatter, bodyText, limit)
		if err == nil {
			return embeddings, nil
		}

		if !errors.Is(err, ai.ErrExceededContextLength) {
			return nil, err
		}

		if limit/2 < mn {
			return nil, err
		}
	}
}

// embedChunks splits the note into chunks of at most limit characters, returning an embedding for each chunk.
func (d *DB) embedChunks(ctx context.Context, frontmatter, body string, limit int) ([][]byte, error) {
	c := chunker.New(
		frontmatter,
		limit,
	)

	for _, b := range chunker.Blocks(body) {
		c.Add(b)
	}

	var (
		chunks     = c.Chunks()
		embeddings = make([][]byte, 0, len(chunks))
	)

	for i, chunk := range chunks {
		serial, err := d.embedChunk(ctx, chunk)
		if err != nil {
			return nil, fmt.Errorf("failed to embed chunk %d/%d: %w", i+1, len(chunks), err)
		}

		// This particular chunk didn't receive an embedding, skip just it.
		if serial == nil {
			continue
		}

		embeddings = append(embeddings, serial)
	}

	return embeddings, nil
}

// embedChunk generates a serialized embedding for a single chunk, returning nil if the model produced no embedding.
func (d *DB) embedChunk(ctx context.Context, c string) ([]byte, error) {
	vec, err := d.client.Embed(ctx, c)
	if err != nil {
		return nil, err
	}

	if len(vec) == 0 {
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
