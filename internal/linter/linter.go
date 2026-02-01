package linter

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/regex"
)

// LintError defines a single linting error.
//
// TODO (jamesl33): Add line number?
type LintError struct {
	Path    string
	Message string
}

// Linter is a struct that contains the logic for linting notes.
type Linter struct{}

// NewLinter creates a new Linter.
func NewLinter() *Linter {
	return &Linter{}
}

// Lint performs linting of notes and returns a slice of linting errors.
func (l *Linter) Lint(ctx context.Context, path string) ([]*LintError, error) {
	ids := make([]string, 0)

	lstr, err := lister.NewLister(
		lister.WithPath(path),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lstr.Many(ctx), hs.Infallible(func(n *note.Note) {
		ids = append(ids, n.Name())
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	entire, err := matcher.Entire("", "", regex.Link.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create entire matcher: %w", err)
	}

	lstr, err = lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(entire),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	errors := make([]*LintError, 0)

	err = iterator.ForEach2(lstr.Many(ctx), func(n *note.Note) error {
		links, err := n.Links()
		if err != nil {
			return fmt.Errorf("failed to get links: %w", err)
		}

		for _, link := range hs.Difference(links, ids) {
			err := LintError{
				Path:    n.Path,
				Message: fmt.Sprintf("Link %q is broken (linkcheck)", link),
			}

			errors = append(errors, &err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	return errors, nil
}
