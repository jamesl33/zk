package note

import (
	"context"

	"github.com/spf13/cobra"
)

// LinksOptions - TODO
type LinksOptions struct {
	// To - TODO
	To bool

	// From - TODO
	From bool
}

// Links - TODO
type Links struct {
	LinksOptions
}

// NewLinks - TODO
func NewLinks() *cobra.Command {
	var links Links

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "links",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return links.Run(cmd.Context(), args) },
	}

	cmd.Flags().BoolVar(
		&links.To,
		"to",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&links.From,
		"from",
		false,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (l *Links) Run(ctx context.Context, args []string) error {
	return nil
}
