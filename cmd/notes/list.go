package notes

import (
	"context"

	"github.com/spf13/cobra"
)

// ListOptions - TODO
type ListOptions struct {
	// Fixed - TODO
	Fixed bool

	// Regex - TODO
	Regex bool

	// Glob - TODO
	Glob bool
}

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
		Args: cobra.MaximumNArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return list.Run(cmd.Context(), args) },
	}

	cmd.Flags().BoolVar(
		&list.Fixed,
		"fixed",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&list.Regex,
		"regex",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&list.Glob,
		"glob",
		false,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (l *List) Run(ctx context.Context, args []string) error {
	return nil
}
