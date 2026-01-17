package note

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesl33/zk/internal/note"
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
func (c *CreateIndex) Run(ctx context.Context, path string) error {
	fm := note.Frontmatter{
		Type:  "index",
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
