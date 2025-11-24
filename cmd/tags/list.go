package tags

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/jamesl33/zk/internal/notes/lister"
	"github.com/spf13/cobra"
)

// ListOptions - TODO
type ListOptions struct{}

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
		RunE: func(cmd *cobra.Command, args []string) error { return list.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run - TODO
func (l *List) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	tags := make(map[string]struct{})

	for n, err := range lister.Many(ctx) {
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		// TODO
		//
		// TODO (jamesl33): Tidy this up.
		for _, tag := range n.Frontmatter.Tags {
			tags[tag] = struct{}{}
		}
	}

	// TODO
	keys := maps.Keys(tags)

	// TODO
	sorted := slices.Sorted(keys)

	// TODO
	compacted := slices.Compact(sorted)

	for _, tag := range compacted {
		fmt.Println(tag)
	}

	return nil
}
