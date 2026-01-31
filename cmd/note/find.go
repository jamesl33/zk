package note

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
	"github.com/spf13/cobra"
)

// FindOptions defines the options for the find command.
type FindOptions struct{}

// Find defines the struct for the find command.
type Find struct {
	FindOptions
}

// NewFind creates a new command for finding semantically similar notes.
func NewFind() *cobra.Command {
	var find Find

	cmd := cobra.Command{
		Short: "Finds semantically similar notes",
		Use:   "find <path>",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return find.Run(cmd.Context(), args[0]) },
	}

	return &cmd
}

// Run finds some related notes.
func (f *Find) Run(ctx context.Context, path string) error {
	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("failed to open note: %w", err)
	}

	err = notes.Find(ctx, n, hs.Infallible(func(n *note.Note) {
		fmt.Println(n.String0())
	}))
	if err != nil {
		return err
	}

	return nil
}
