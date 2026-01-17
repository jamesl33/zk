package note

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesl33/zk/internal/note"
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
func (c *CreateLiterature) Run(ctx context.Context, path string) error {
	fm := note.Frontmatter{
		Type:  "literature",
		Title: c.Title,
		Date:  time.Now().Format("2006-01-02"),
		Tags:  make([]string, 0),
	}

	n := note.Note{
		Path:        note.Path(path),
		Frontmatter: fm,
	}

	err := n.Write()
	if err != nil {
		return fmt.Errorf("failed to write note: %w", err)
	}

	fmt.Printf("%s\n", n.Path)

	return nil
}
