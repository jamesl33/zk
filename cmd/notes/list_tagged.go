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

// ListTaggedOptions defines the options for the tagged command.
type ListTaggedOptions struct {
	// With defines a list of tags which must be present.
	With []string

	// Without defines a list of tags which must not be present.
	Without []string
}

// ListTagged defines the struct for the tagged command.
type ListTagged struct {
	ListTaggedOptions
}

// NewListTagged creates a new command for listing notes by tag.
func NewListTagged() *cobra.Command {
	var tagged ListTagged

	cmd := cobra.Command{
		Short: "List notes by tag",
		Use:   "tagged",
		RunE:  func(cmd *cobra.Command, args []string) error { return tagged.Run(cmd.Context(), args) },
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
		fmt.Println(n.String())
	}))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	return nil
}
