package index

import (
	"context"
	"database/sql"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strings"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes/lister"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
)

// TODO (jamesl33): Create an 'internal/database' package.
func init() {
	sqlite_vec.Auto()
}

// CreateOptions - TODO
type CreateOptions struct{}

// Create - TODO
type Create struct {
	CreateOptions
}

// NewCreate - TODO
func NewCreate() *cobra.Command {
	var create Create

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "create",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return create.Run(cmd.Context()) },
	}

	return &cmd
}

// Run - TODO
func (c *Create) Run(ctx context.Context) error {
	db, err := c.open(ctx)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}
	defer db.Close()

	err = c.populate(ctx, db)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// open - TODO
func (c *Create) open(ctx context.Context) (*sql.DB, error) {
	// path - TODO
	const path = "zk.sqlite3"

	db, err := sql.Open("sqlite3", filepath.Join(os.TempDir(), path))
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// create - TODO
	const create = `CREATE table IF NOT EXISTS notes (name text, checksum integer, embedding blob not null);`

	_, err = db.ExecContext(ctx, create)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return db, nil
}

// populate - TODO
func (c *Create) populate(ctx context.Context, db *sql.DB) error {
	lister, err := lister.NewLister(
		lister.WithPath("."),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	for n, err := range lister.Many(ctx) {
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		err = c.insert(ctx, db, n)
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}
	}

	return nil
}

// insert - TODO
func (c *Create) insert(ctx context.Context, db *sql.DB, n *note.Note) error {
	// TODO
	hasher := crc32.NewIEEE()

	_, err := io.Copy(hasher, strings.NewReader(n.Body))
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	embedding, err := c.embed(ctx, n)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if embedding == nil {
		return nil
	}

	serial, err := sqlite_vec.SerializeFloat32(embedding)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// insert - TODO
	const insert = `INSERT INTO notes VALUES (?, ?, ?);`

	_, err = db.ExecContext(ctx, insert, n.Name(), hasher.Sum32(), serial)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// embed - TODO
//
// TODO (jamesl33): Create a vector search package.
func (c *Create) embed(ctx context.Context, n *note.Note) ([]float32, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// TODO (jamesl33): Handle the 2k context window.
	req := api.EmbedRequest{
		Model: "embeddinggemma:300m",
		Input: n.Body,
	}

	resp, err := client.Embed(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if len(resp.Embeddings) != 1 {
		return nil, nil
	}

	return resp.Embeddings[0], nil
}
