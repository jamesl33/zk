package index

import (
	"context"

	"github.com/spf13/cobra"
)

// UpdateOptions - TODO
type UpdateOptions struct{}

// Update - TODO
type Update struct {
	UpdateOptions
}

// NewUpdate - TODO
func NewUpdate() *cobra.Command {
	var update Update

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "update",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return update.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run - TODO
func (u *Update) Run(ctx context.Context, args []string) error {
	return nil
}
