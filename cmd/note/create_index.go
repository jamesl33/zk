package note

import (
	"context"

	"github.com/spf13/cobra"
)

// CreateIndexOptions defines the options for the index command.
type CreateIndexOptions struct {
	// Title is the title for the note (e.g. the title of a book/article).
	Title string
}

// CreateIndex defines the struct for the index command.
type CreateIndex struct {
	CreateIndexOptions
}

// NewCreateIndex creates a new command for creating an 'index' note.
func NewCreateIndex() *cobra.Command {
	var index CreateIndex

	cmd := cobra.Command{
		Short: "Create a new 'index' note",
		Use:   "index <directory>",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return index.Run(cmd.Context(), args[0]) },
	}

	cmd.Flags().StringVar(
		&index.Title,
		"title",
		"Untitled",
		"The title for the note (e.g. the title of a book/article)",
	)

	return &cmd
}

// Run creates a new index note.
func (c *CreateIndex) Run(_ context.Context, path string) error {
	return create("index", c.Title, path)
}
