package note

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// CreateIndexOptions - TODO
type CreateIndexOptions struct {
	// Title - TODO
	Title string
}

// CreateIndex - TODO
type CreateIndex struct {
	CreateIndexOptions
}

// NewCreateIndex - TODO
func NewCreateIndex() *cobra.Command {
	var index CreateIndex

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "index",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return index.Run(cmd.Context(), args[0]) },
	}

	cmd.Flags().StringVar(
		&index.Title,
		"title",
		"Untitled",
		// TODO
		"",
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
