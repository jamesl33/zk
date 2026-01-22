package notes

import "github.com/spf13/cobra"

// NewNotes creates a new command for interacting with notes.
func NewNotes() *cobra.Command {
	cmd := cobra.Command{
		Short: "Interact with all the notes in the Zettelkasten",
		Use:   "notes",
	}

	cmd.AddCommand(
		NewFind(),
		NewLint(),
		NewList(),
		NewPick(),
		NewSearch(),
	)

	return &cmd
}
