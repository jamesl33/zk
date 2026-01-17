package note

import "github.com/spf13/cobra"

// NewNote creates a new command for interacting with a single note.
func NewNote() *cobra.Command {
	cmd := cobra.Command{
		Short: "Interact and manipulate a single note",
		Use:   "note",
	}

	cmd.AddCommand(
		NewCreate(),
		NewUpdate(),
		NewLinks(),
		NewFind(),
		NewSummarize(),
	)

	return &cmd
}
