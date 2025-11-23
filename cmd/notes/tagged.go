package notes

import (
	"context"

	"github.com/spf13/cobra"
)

// TaggedOptions - TODO
type TaggedOptions struct {
	// With - TODO
	With bool

	// Without - TODO
	Without bool
}

// Tagged - TODO
type Tagged struct {
	TaggedOptions
}

// NewTagged - TODO
func NewTagged() *cobra.Command {
	var tagged Tagged

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "tagged",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return tagged.Run(cmd.Context(), args) },
	}

	cmd.Flags().BoolVar(
		&tagged.With,
		"with",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&tagged.Without,
		"without",
		false,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (t *Tagged) Run(ctx context.Context, args []string) error {
	return nil
}
