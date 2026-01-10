package note

import "github.com/spf13/cobra"

// NewNote - TODO
func NewNote() *cobra.Command {
	cmd := cobra.Command{
		// TODO
		Short: "Interact and manipulate a single note",
		// TODO
		Use: "note",
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
