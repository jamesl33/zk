package note

import "github.com/spf13/cobra"

// NewNote - TODO
func NewNote() *cobra.Command {
	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "note",
	}

	cmd.AddCommand(
		NewCreate(),
		NewUpdate(),
		NewLinks(),
		NewSummarize(),
	)

	return &cmd
}
