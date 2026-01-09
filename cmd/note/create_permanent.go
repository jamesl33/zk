package note

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// CreatePermanentOptions - TODO
type CreatePermanentOptions struct {
	// Title - TODO
	Title string
}

// CreatePermanent - TODO
type CreatePermanent struct {
	CreatePermanentOptions
}

// NewCreatePermanent - TODO
func NewCreatePermanent() *cobra.Command {
	var permanent CreatePermanent

	cmd := cobra.Command{
		// TODO
		Short: "Create a new 'permanent' note",
		// TODO
		Use: "permanent directory",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return permanent.Run(cmd.Context(), args[0]) },
	}

	cmd.Flags().StringVar(
		&permanent.Title,
		"title",
		"Untitled",
		"The title for the note (e.g. the title of a book/article)",
	)

	return &cmd
}

// Run creates a new permanent note.
func (c *CreatePermanent) Run(ctx context.Context, path string) error {
	fm := note.Frontmatter{
		Type:  "permanent",
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
