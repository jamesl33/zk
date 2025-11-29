package index

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jamesl33/zk/internal/database/vector"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/spf13/cobra"
)

// IndexOptions - TODO
type IndexOptions struct{}

// Index - TODO
type Index struct {
	IndexOptions
}

// NewIndex - TODO
func NewIndex() *cobra.Command {
	var index Index

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "index",
		// TODO
		RunE: func(cmd *cobra.Command, _ []string) error { return index.Run(cmd.Context()) },
	}

	return &cmd
}

// Run - TODO
func (i *Index) Run(ctx context.Context) error {
	db, err := vector.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}
	defer db.Close()

	err = i.populate(ctx, db)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// populate - TODO
func (i *Index) populate(ctx context.Context, db *vector.DB) error {
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
