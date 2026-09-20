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

// index is the subset of vector.DB's behavior that Finder depends on, allowing it to be
// substituted with a mock/fake (backed by a mock AI 'Client') in tests.
type index interface {
	Upsert(ctx context.Context, n *note.Note) error
	Find(ctx context.Context, n *note.Note) ([]*note.Note, error)
	Close() error
}

var _ index = (*vector.DB)(nil)

// Finder finds notes which are semantically similar to a given note, using a vector index.
type Finder struct {
	db index
}

// NewFinder creates a Finder backed by the vector index at '.zk/zk.sqlite3'.
func NewFinder(ctx context.Context) (*Finder, error) {
	db, err := vector.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return &Finder{db: db}, nil
}

// Close releases resources held by the Finder.
func (f *Finder) Close() error {
	return f.db.Close()
}

// Find finds notes which are semantically similar to the given note. It handles populating the
// index and finding related notes.
func (f *Finder) Find(ctx context.Context, n *note.Note, fn func(n *note.Note) error) error {
	err := f.populate(ctx)
	if err != nil {
		return fmt.Errorf("failed to populate database: %w", err)
	}

	notes, err := f.db.Find(ctx, n)
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
func (f *Finder) populate(ctx context.Context) error {
	lister, err := lister.NewLister(
		lister.WithPath("."),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), func(n *note.Note) error {
		return f.db.Upsert(ctx, n)
	})
	if err != nil {
		return fmt.Errorf("failed to upsert embeddings: %w", err)
	}

	return nil
}

// Find finds notes which are semantically similar to the given note. It handles creating the
// database, populating it, and finding related notes.
func Find(ctx context.Context, n *note.Note, fn func(n *note.Note) error) error {
	finder, err := NewFinder(ctx)
	if err != nil {
		return err
	}
	defer finder.Close()

	return finder.Find(ctx, n, fn)
}
