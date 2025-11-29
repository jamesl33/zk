package tags

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/note"
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
		Args: cobra.MaximumNArgs(1),
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

	all := make(map[string]struct{})

	cp := func(tags []string) {
		iterator.ForEach(slices.Values(tags), func(tag string) { all[tag] = struct{}{} })
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		cp(n.Frontmatter.Tags)
	}))
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	keys := maps.Keys(all)

	// TODO
	sorted := slices.Sorted(keys)

	// TODO
	compacted := slices.Compact(sorted)

	for _, tag := range compacted {
		fmt.Println(tag)
	}

	return nil
}
