package tags

import (
	"context"

	"github.com/spf13/cobra"
)

// GenerateOptions - TODO
type GenerateOptions struct{}

// Generate - TODO
type Generate struct {
	GenerateOptions
}

// NewGenerate - TODO
func NewGenerate() *cobra.Command {
	var generate Generate

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "generate",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return generate.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run - TODO
func (g *Generate) Run(ctx context.Context, args []string) error {
	return nil
}
