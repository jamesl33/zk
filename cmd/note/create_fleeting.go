package note

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// CreateFleetingOptions - TODO
type CreateFleetingOptions struct {
	// Title - TODO
	Title string
}

// CreateFleeting - TODO
type CreateFleeting struct {
	CreateFleetingOptions
}

// NewCreateFleeting - TODO
func NewCreateFleeting() *cobra.Command {
	var fleeting CreateFleeting

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "fleeting",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return fleeting.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringVar(
		&fleeting.Title,
		"title",
		"Untitled",
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (c *CreateFleeting) Run(ctx context.Context, args []string) error {
	fm := note.Frontmatter{
		Type:  "fleeting",
		Title: c.Title,
		Date:  time.Now().Format("2006-01-02"),
		Tags:  make([]string, 0),
	}

	n := note.Note{
		Path:        note.Path("0 Inbox"),
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

	fmt.Println(n.Body)

	return nil
}
