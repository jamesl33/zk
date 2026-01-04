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

// ListOptions - TODO
//
// TODO (jamesl33): Add support for case-insensitive listing.
// TODO (jamesl33): Don't use boolean flags, use string flags then argument can be for the sub-directory.
type ListOptions struct {
	// Fixed - TODO
	Fixed string

	// Glob - TODO
	Glob string

	// Regex - TODO
	Regex string
}

// List - TODO
type List struct {
	ListOptions
}

// NewList - TODO
func NewList() *cobra.Command {
	var list List

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "list",
		// TODO
		Args: cobra.MaximumNArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return list.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringVar(
		&list.Fixed,
		"fixed",
		"",
		// TODO
		"",
	)

	cmd.Flags().StringVar(
		&list.Glob,
		"glob",
		"",
		// TODO
		"",
	)

	cmd.Flags().StringVar(
		&list.Regex,
		"regex",
		"",
		// TODO
		"",
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

	title, err := matcher.Title(l.Fixed, l.Glob, l.Regex)
	if err != nil {
		return fmt.Errorf("failed to create matcher: %w", err)
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(title),
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
