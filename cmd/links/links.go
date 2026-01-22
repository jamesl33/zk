package links

import "github.com/spf13/cobra"

// NewLinks creates a new command for interacting with note links.
func NewLinks() *cobra.Command {
	cmd := cobra.Command{
		Short: "Interact and manipulate note links",
		Use:   "links",
	}

	cmd.AddCommand(
		NewList(),
		NewRewrite(),
	)

	return &cmd
}
