package notes

import (
	"context"

	"github.com/spf13/cobra"
)

// FindOptions - TODO
type FindOptions struct {
	// Similar - TODO
	Similar bool

	// Dissimilar - TODO
	Dissimilar bool
}

// Find - TODO
type Find struct {
	FindOptions
}

// NewFind - TODO
func NewFind() *cobra.Command {
	var find Find

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "find",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return find.Run(cmd.Context(), args) },
	}

	cmd.Flags().BoolVar(
		&find.Similar,
		"similar",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&find.Dissimilar,
		"dissimilar",
		false,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (f *Find) Run(ctx context.Context, args []string) error {
	return nil
}
