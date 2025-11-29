package tags

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes/lister"
	"github.com/jamesl33/zk/internal/notes/matcher"
	"github.com/spf13/cobra"
)

// DeleteOptions - TODO
type DeleteOptions struct{}

// Delete - TODO
type Delete struct {
	DeleteOptions
}

// NewDelete - TODO
func NewDelete() *cobra.Command {
	var del Delete

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "delete",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return del.Run(cmd.Context(), args[0]) },
	}

	return &cmd
}

// Run - TODO
func (d *Delete) Run(ctx context.Context, remove string) error {
	tags, err := matcher.Tags([]string{remove}, nil)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(tags),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = iterator.ForEach2(lister.Many(ctx), func(n *note.Note) error {
		return d.update(n, remove)
	})
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// update - TODO
func (d *Delete) update(n *note.Note, remove string) error {
	n.Frontmatter.Tags = hs.Filter(n.Frontmatter.Tags, func(tag string) bool { return tag != remove })

	// TODO
	err := n.Write()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}
