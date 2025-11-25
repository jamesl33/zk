package note

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// CreateBibliographicOptions - TODO
type CreateBibliographicOptions struct {
	// Title - TODO
	Title string
}

// CreateBibliographic - TODO
type CreateBibliographic struct {
	CreateBibliographicOptions
}

// NewCreateBibliographic - TODO
func NewCreateBibliographic() *cobra.Command {
	var bibliographic CreateBibliographic

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "bibliographic",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return bibliographic.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringVar(
		&bibliographic.Title,
		"title",
		"Untitled",
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
//
// TODO (jamesl33): Auto-open the created note.
func (c *CreateBibliographic) Run(ctx context.Context, args []string) error {
	fm := note.Frontmatter{
		Type:  "bibliographic",
		Title: c.Title,
		Date:  time.Now().Format("2006-01-02"),
		Tags:  make([]string, 0),
	}

	n := note.Note{
		Path:        note.Path("5 Bibliography"),
		Frontmatter: fm,
	}

	err := n.Write()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}
