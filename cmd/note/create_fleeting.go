package note

import (
	"context"
	"fmt"
	"time"

	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// CreateFleetingOptions defines the options for the fleeting command.
type CreateFleetingOptions struct {
	// Title is the title for the note.
	Title string
}

// CreateFleeting defines the struct for the fleeting command.
type CreateFleeting struct {
	CreateFleetingOptions
}

// NewCreateFleeting creates a new command for creating a 'fleeting' note.
func NewCreateFleeting() *cobra.Command {
	var fleeting CreateFleeting

	cmd := cobra.Command{
		Short: "Create a new 'fleeting' note",
		Use:   "fleeting",
		RunE:  func(cmd *cobra.Command, _ []string) error { return fleeting.Run(cmd.Context()) },
	}

	cmd.Flags().StringVar(
		&fleeting.Title,
		"title",
		"Untitled",
		"The title for the note",
	)

	return &cmd
}

// Run creates a new fleeting note.
func (c *CreateFleeting) Run(ctx context.Context) error {
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
		return fmt.Errorf("failed to write note: %w", err)
	}

	fmt.Printf("%s\n", n.Path)

	return nil
}
