package tags

import (
	"context"

	"github.com/spf13/cobra"
)

// ListOptions - TODO
type ListOptions struct{}

// List - TODO
type List struct {
	ListOptions
}

// NewList - TODO
func NewList() *cobra.Command {
	var list List

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "list",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return list.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run - TODO
func (l *List) Run(ctx context.Context, args []string) error {
	return nil
}
