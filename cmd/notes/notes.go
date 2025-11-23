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
		// TODO
		NewList(),
		// TODO
		NewSearch(),
		// TODO
		NewFind(),
		// TODO
		NewTagged(),
		// TODO
		NewPick(),
	)

	return &cmd
}
