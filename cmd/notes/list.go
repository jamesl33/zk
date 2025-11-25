package notes

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/notes/lister"
	"github.com/jamesl33/zk/internal/notes/matcher"
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

	return &cmd
}

// Run - TODO
func (l *List) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	title, err := matcher.Title(l.Fixed, l.Glob, l.Regex)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(title),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	for n, err := range lister.Many(ctx) {
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		fmt.Println(n.String0())
	}

	return nil
}
