package note

import (
	"github.com/spf13/cobra"
)

// NewCreate - TODO
func NewCreate() *cobra.Command {
	cmd := cobra.Command{
		// TODO
		Short: "Create a new note",
		// TODO
		Use: "create",
	}

	cmd.AddCommand(
		NewCreateBibliographic(),
		NewCreatePermanent(),
		NewCreateFleeting(),
		NewCreateIndex(),
		NewCreateLiterature(),
	)

	return &cmd
}
