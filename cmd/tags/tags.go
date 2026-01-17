package tags

import "github.com/spf13/cobra"

// NewTags creates a new command for interacting with note tags.
func NewTags() *cobra.Command {
	cmd := cobra.Command{
		Short: "Interact and manipulate note tags",
		Use:   "tags",
	}

	cmd.AddCommand(
		NewGenerate(),
		NewList(),
		NewDelete(),
	)

	return &cmd
}
