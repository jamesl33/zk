package note

import (
	"context"

	"github.com/spf13/cobra"
)

// CreateOptions - TODO
type CreateOptions struct {
	// Template - TODO
	Template string
}

// Create - TODO
type Create struct {
	CreateOptions
}

// NewCreate - TODO
func NewCreate() *cobra.Command {
	var create Create

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "create",
		// TODO
		Args: cobra.ExactArgs(2),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return create.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringVar(
		&create.Template,
		"template",
		"",
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (c *Create) Run(ctx context.Context, args []string) error {
	return nil
}
