package index

import (
	"context"

	"github.com/spf13/cobra"
)

// DeleteOptions - TODO
type DeleteOptions struct{}

// Delete - TODO
type Delete struct {
	DeleteOptions
}

// NewDelete - TODO
func NewDelete() *cobra.Command {
	var del Delete

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "delete",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return del.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run - TODO
func (u *Delete) Run(ctx context.Context, args []string) error {
	return nil
}
