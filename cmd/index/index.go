package index

import "github.com/spf13/cobra"

// NewIndex - TODO
func NewIndex() *cobra.Command {
	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "index",
	}

	cmd.AddCommand(
		// TODO
		NewCreate(),
		// TODO
		NewUpdate(),
		// TODO
		NewDelete(),
	)

	return &cmd
}
