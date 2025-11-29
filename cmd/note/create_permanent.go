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
		Short: "",
		// TODO
		Use: "permanent",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return permanent.Run(cmd.Context(), args[0]) },
	}

	cmd.Flags().StringVar(
		&permanent.Title,
		"title",
		"Untitled",
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
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
		return fmt.Errorf("%w", err) // TODO
	}

	fmt.Printf("%s\n", n.Path)

	return nil
}
