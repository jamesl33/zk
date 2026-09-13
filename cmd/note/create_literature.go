package note

import (
	"context"

	"github.com/spf13/cobra"
)

// CreateLiteratureOptions defines the options for the literature command.
type CreateLiteratureOptions struct {
	// Title is the title for the note (e.g. the title of a book/article).
	Title string
}

// CreateLiterature defines the struct for the literature command.
type CreateLiterature struct {
	CreateLiteratureOptions
}

// NewCreateLiterature creates a new command for creating a 'literature' note.
func NewCreateLiterature() *cobra.Command {
	var literature CreateLiterature

	cmd := cobra.Command{
		Short: "Create a new 'literature' note",
		Use:   "literature <directory>",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return literature.Run(cmd.Context(), args[0]) },
	}

	cmd.Flags().StringVar(
		&literature.Title,
		"title",
		"Untitled",
		"The title for the note (e.g. the title of a book/article)",
	)

	return &cmd
}

// Run creates a new literature note.
func (c *CreateLiterature) Run(_ context.Context, path string) error {
	return create("literature", c.Title, path)
}
