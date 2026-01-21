package links

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/links"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// RewriteOptions defines the options for the rewrite command.
type RewriteOptions struct{}

// Rewrite defines the struct for the rewrite command.
type Rewrite struct {
	RewriteOptions
}

// NewRewrite creates a new command for rewriting note links.
func NewRewrite() *cobra.Command {
	var rewrite Rewrite

	cmd := cobra.Command{
		Short: "Rewrite note links to use the note title",
		Use:   "rewrite [directory | path]",
		Args:  cobra.MaximumNArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return rewrite.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run the command to find linked notes.
func (r *Rewrite) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	l, err := lister.NewLister(
		lister.WithPath(path),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(l.Many(ctx), func(n *note.Note) error {
		return r.rewrite(ctx, n)
	})
	if err != nil {
		return fmt.Errorf("failed to rewrite links: %w", err)
	}

	return nil
}

// rewrite the links in the given note.
func (r *Rewrite) rewrite(ctx context.Context, n *note.Note) error {
	err := links.Rewrite(ctx, n)
	if err != nil {
		return fmt.Errorf("failed to rewrite links: %w", err)
	}

	err = n.Write()
	if err != nil {
		return fmt.Errorf("failed to write note: %w", err)
	}

	return nil
}
