package notes

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/regex"
	"github.com/spf13/cobra"
)

// LintOptions defines the options for the lint command.
type LintOptions struct{}

// Lint defines the struct for the lint command.
type Lint struct {
	LintOptions
}

// NewLint creates a new command for linting notes.
func NewLint() *cobra.Command {
	var lint Lint

	cmd := cobra.Command{
		Short: "Lints notes",
		Use:   "lint [directory]",
		Args:  cobra.MaximumNArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return lint.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run lints the notes, printing warnings/errors.
func (l *Lint) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	ids := make([]string, 0)

	lstr, err := lister.NewLister(
		lister.WithPath(path),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lstr.Many(ctx), hs.Infallible(func(n *note.Note) {
		ids = append(ids, n.Name())
	}))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	entire, err := matcher.Entire("", "", regex.Link.String())
	if err != nil {
		return fmt.Errorf("failed to create entire matcher: %w", err)
	}

	lstr, err = lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(entire),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	// TODO (jamesl33): Mimic 'golangci-lint' more closely, by reporting the line/offset as well?
	// TODO (jamesl33): Refactor this to allow modular linting.
	err = iterator.ForEach2(lstr.Many(ctx), hs.Infallible(func(n *note.Note) {
		for _, link := range hs.Difference(n.Links(), ids) {
			fmt.Printf("%q: Link %q is broken (linkcheck)\n", n.Path, link)
		}
	}))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	return nil
}
