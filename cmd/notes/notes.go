package notes

import "github.com/spf13/cobra"

// NewNotes - TODO
func NewNotes() *cobra.Command {
	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "notes",
	}

	cmd.AddCommand(
		NewList(),
		NewSearch(),
		NewFind(),
		NewTagged(),
		NewPick(),
	)

	return &cmd
}
