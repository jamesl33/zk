package tags

import "github.com/spf13/cobra"

// NewTags - TODO
func NewTags() *cobra.Command {
	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "tags",
	}

	cmd.AddCommand(
		// TODO
		NewGenerate(),
		// TODO
		NewList(),
		// TODO
		NewDelete(),
	)

	return &cmd
}
