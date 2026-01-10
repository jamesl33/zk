package note

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// CreateLiteratureOptions - TODO
type CreateLiteratureOptions struct {
	// Title - TODO
	Title string
}

// CreateLiterature - TODO
type CreateLiterature struct {
	CreateLiteratureOptions
}

// NewCreateLiterature - TODO
func NewCreateLiterature() *cobra.Command {
	var literature CreateLiterature

	cmd := cobra.Command{
		// TODO
		Short: "Create a new 'literature' note",
		// TODO
		Use: "literature <directory>",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return literature.Run(cmd.Context(), args[0]) },
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
