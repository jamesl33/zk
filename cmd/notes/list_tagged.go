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

// ListTaggedOptions - TODO
type ListTaggedOptions struct {
	// With - TODO
	With []string

	// Without - TODO
	Without []string
}

// ListTagged - TODO
type ListTagged struct {
	ListTaggedOptions
}

// NewListTagged - TODO
func NewListTagged() *cobra.Command {
	var tagged ListTagged

	cmd := cobra.Command{
		// TODO
		Short: "List notes by tag",
		// TODO
		Use: "tagged",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return tagged.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringArrayVar(
		&tagged.With,
		"with",
		nil,
		"Include notes which have the provided tag",
	)

	cmd.Flags().StringArrayVar(
		&tagged.Without,
		"without",
		nil,
		"Exclude notes which have the provided tag",
	)

	return &cmd
}

// Run lists tagged notes.
func (l *ListTagged) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	tags, err := matcher.Tags(l.With, l.Without)
	if err != nil {
		return fmt.Errorf("failed to create matcher: %w", err)
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(tags),
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
