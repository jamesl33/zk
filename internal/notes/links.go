package notes

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
)

// LinkedFrom finds notes which are linked from the given note and calls the function for each note.
func LinkedFrom(ctx context.Context, n *note.Note, fn func(n *note.Note)) error {
	matchers := hs.Map(n.Links(), func(n string) matcher.Matcher { return matcher.Name(n) })

	// Must check for no matchers, as the default is to list all
	if len(matchers) == 0 {
		return nil
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(fn))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	return nil
}

// LinkedTo finds notes which link to the given note and calls the function for each note.
func LinkedTo(ctx context.Context, n *note.Note, fn func(n *note.Note)) error {
	var (
		// name of the note, escaped for use in regular expressions
		name = regexp.QuoteMeta(n.Name())

		// pattern which matches links to this note
		pattern = fmt.Sprintf(`\[\[%s(\|.*?)?\]\]`, name)
	)

	matcher, err := matcher.Body("", "", pattern)
	if err != nil {
		return fmt.Errorf("failed to create matcher: %w", err)
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(fn))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	return nil
}
