package note

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vector"
	"github.com/spf13/cobra"
)

// FindOptions - TODO
type FindOptions struct{}

// Find - TODO
type Find struct {
	FindOptions
}

// NewFind - TODO
func NewFind() *cobra.Command {
	var find Find

	cmd := cobra.Command{
		// TODO
		Short: "Finds semantically similar notes",
		// TODO
		Use: "find <path>",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return find.Run(cmd.Context(), args[0]) },
	}

	return &cmd
}

// Run finds some related notes.
func (f *Find) Run(ctx context.Context, path string) error {
	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("failed to open note: %w", err)
	}

	db, err := vector.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	err = f.populate(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to populate database: %w", err)
	}

	notes, err := db.Find(ctx, n)
	if err != nil {
		return fmt.Errorf("failed to find related notes: %w", err)
	}

	err = iterator.ForEach2(notes, hs.Infallible(func(n *note.Note) {
		fmt.Println(n.String0())
	}))
	if err != nil {
		return fmt.Errorf("failed to list related notes: %w", err)
	}

	return nil
}

// populate the index by updating embeddings for notes that have been updated.
func (f *Find) populate(ctx context.Context, db *vector.DB) error {
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
