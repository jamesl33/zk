package notes

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vector"
)

// Find finds notes which are semantically similar to the given note. It handles creating the database,
// populating it, and finding related notes.
//
// TODO (jamesl33): Extract a struct out of this, to enable unit testing with a mock AI 'Client'.
func Find(ctx context.Context, n *note.Note, fn func(n *note.Note) error) error {
	db, err := vector.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	err = populate(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to populate database: %w", err)
	}

	notes, err := db.Find(ctx, n)
	if err != nil {
		return fmt.Errorf("failed to find related notes: %w", err)
	}

	for _, n := range notes {
		if err := fn(n); err != nil {
			return err
		}
	}

	return nil
}

// populate the index by updating embeddings for notes that have been updated.
func populate(ctx context.Context, db *vector.DB) error {
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
