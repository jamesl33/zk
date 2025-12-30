package index

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vector"
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

// Run index creation.
//
// TODO (jamesl33): Populate summaries?
func (i *Index) Run(ctx context.Context) error {
	db, err := vector.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	err = i.populate(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to populate database: %w", err)
	}

	return nil
}

// populate the index.
func (i *Index) populate(ctx context.Context, db *vector.DB) error {
	lister, err := lister.NewLister(
		lister.WithPath("."),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), func(n *note.Note) error {
		return db.Upsert(ctx, n)
	})
	if err != nil {
		return fmt.Errorf("failed to upsert embeddings: %w", err)
	}

	return nil
}
