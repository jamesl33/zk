package memory

import (
	"context"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
	"github.com/jamesl33/zk/internal/vector"
)

// Finder finds the notes closest to a query, no further than maxDistance away. It runs in the vault.
type Finder func(ctx context.Context, query string, limit int, maxDistance float32) ([]*note.Note, error)

// Find is the Finder backed by the vault's vector index.
func Find(ctx context.Context, query string, limit int, maxDistance float32) ([]*note.Note, error) {
	var n note.Note

	n.SetBody(query)

	var found []*note.Note

	err := notes.Find(ctx, &n, func(n *note.Note) error {
		found = append(found, n)

		return nil
	}, vector.WithLimit(limit), vector.WithMaxDistance(maxDistance))

	return found, err
}
