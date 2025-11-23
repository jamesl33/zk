package notes

import (
	"context"

	"github.com/spf13/cobra"
)

// SearchOptions - TODO
type SearchOptions struct {
	// Fixed - TODO
	Fixed bool

	// Regex - TODO
	Regex bool
}

// Search - TODO
type Search struct {
	SearchOptions
}

// NewSearch - TODO
func NewSearch() *cobra.Command {
	var search Search

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "search",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return search.Run(cmd.Context(), args) },
	}

	cmd.Flags().BoolVar(
		&search.Fixed,
		"fixed",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&search.Regex,
		"regex",
		false,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (s *Search) Run(ctx context.Context, args []string) error {
	return nil
}
