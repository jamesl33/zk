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
		Short: "",
		// TODO
		Use: "literature",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return literature.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringVar(
		&literature.Title,
		"title",
		"Untitled",
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (c *CreateLiterature) Run(ctx context.Context, args []string) error {
	fm := note.Frontmatter{
		Type:  "literature",
		Title: c.Title,
		Date:  time.Now().Format("2006-01-02"),
		Tags:  make([]string, 0),
	}

	n := note.Note{
		Path:        note.Path(args[0]),
		Frontmatter: fm,
	}

	err := n.Write()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = n.Edit(ctx)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}
