package notes

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
)

// Search executes a search for notes in a given path with a provided matcher, calling the function for each found note.
func Search(ctx context.Context, path string, matcher matcher.Matcher, fn func(n *note.Note)) error {
	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(matcher),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(fn))
	if err != nil {
		return fmt.Errorf("failed to search notes: %w", err)
	}

	return nil
}
