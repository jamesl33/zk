package links

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

// ListOptions defines the options for the list command.
type ListOptions struct {
	// Broken means only notes with broken links will be shown.
	Broken bool
}

// List defines the struct for the list command.
//
// TODO (jamesl33): Turn this into a "lint" command, and just have the ability to lint notes (then it can use LSP diagnostics as well).
type List struct {
	ListOptions
}

// NewList creates a new command for listing notes with links.
func NewList() *cobra.Command {
	var list List

	cmd := cobra.Command{
		Short: "Lists notes which contain links to other notes",
		Use:   "list [directory]",
		Args:  cobra.MaximumNArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return list.Run(cmd.Context(), args) },
	}

	cmd.Flags().BoolVar(
		&list.Broken,
		"broken",
		false,
		"Filter notes by those which have broken links",
	)

	return &cmd
}

// Run lists notes with matching titles.
func (l *List) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	ids := make([]string, 0)

	lstr, err := lister.NewLister(
		lister.WithPath(path),
	)

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

	err = iterator.ForEach2(lstr.Many(ctx), hs.Infallible(func(n *note.Note) {
		if l.Broken && len(hs.Difference(n.Links(), ids)) == 0 {
			return
		}

		fmt.Println(n.String0())
	}))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	return nil
}
