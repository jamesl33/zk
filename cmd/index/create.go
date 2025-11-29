package index

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/jamesl33/zk/internal/database/vector"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes/lister"
	_ "github.com/mattn/go-sqlite3"
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
	db, err := vector.New(ctx, filepath.Join(os.TempDir(), "zk.sqlite3"))
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

// populate - TODO
func (c *Create) populate(ctx context.Context, db *vector.DB) error {
	lister, err := lister.NewLister(
		lister.WithPath("."),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = iterator.ForEach2(lister.Many(ctx), func(n *note.Note) error {
		return db.Upsert(ctx, n)
	})
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}
