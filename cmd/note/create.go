package note

import (
	"github.com/spf13/cobra"
)

// NewCreate creates a new command for creating notes.
func NewCreate() *cobra.Command {
	cmd := cobra.Command{
		Short: "Create a new note",
		Use:   "create",
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
