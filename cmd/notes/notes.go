package notes

import "github.com/spf13/cobra"

// NewNotes - TODO
func NewNotes() *cobra.Command {
	cmd := cobra.Command{
		// TODO
		Short: "Interact with all the notes in the Zettelkasten",
		// TODO
		Use: "notes",
	}

	cmd.AddCommand(
		NewList(),
		NewSearch(),
		NewFind(),
		NewPick(),
	)

	return &cmd
}
