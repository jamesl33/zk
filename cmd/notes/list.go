package notes

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// ListOptions defines the options for the list command.
//
// TODO (jamesl33): Add support for case-insensitive listing.
type ListOptions struct {
	// Fixed filters notes by title using a case-sensitive fixed-string search.
	Fixed string

	// Glob filters notes by title using a case-sensitive glob pattern.
	Glob string

	// Regex filters notes by title using a regular expression (RE2).
	Regex string
}

// List defines the struct for the list command.
type List struct {
	ListOptions
}

// NewList creates a new command for listing notes.
func NewList() *cobra.Command {
	var list List

	cmd := cobra.Command{
		Short: "List notes and display useful information/metadata (e.g. tags)",
		Use:   "list [directory]",
		Args:  cobra.MaximumNArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return list.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringVar(
		&list.Fixed,
		"fixed",
		"",
		"Filter notes by title using a case-sensitive fixed-string search",
	)

	cmd.Flags().StringVar(
		&list.Glob,
		"glob",
		"",
		"Filter notes by title using a case-sensitive glob pattern",
	)

	cmd.Flags().StringVar(
		&list.Regex,
		"regex",
		"",
		"Filter notes by title using a regular expression (RE2)",
	)

	cmd.AddCommand(
		NewListTagged(),
	)

	return &cmd
}

// Run lists notes with matching titles.
func (l *List) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	pm, err := matcher.Path(l.Fixed, l.Glob, l.Regex)
	if err != nil {
		return fmt.Errorf("failed to create path matcher: %w", err)
	}

	title, err := matcher.Title(l.Fixed, l.Glob, l.Regex)
	if err != nil {
		return fmt.Errorf("failed to create title matcher: %w", err)
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(matcher.Or(pm, title)),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		fmt.Println(n.String0())
	}))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	return nil
}
