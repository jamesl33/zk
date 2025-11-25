package note

import (
	"github.com/spf13/cobra"
)

// NewCreate - TODO
func NewCreate() *cobra.Command {
	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "create",
	}

	cmd.AddCommand(
		// TODO
		NewCreateBibliographic(),
		// TODO
		NewCreatePermanent(),
		// TODO
		NewCreateFleeting(),
		// TODO
		NewCreateIndex(),
		// TODO
		NewCreateLiterature(),
	)

	return &cmd
}
