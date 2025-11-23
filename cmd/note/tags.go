package note

import (
	"context"

	"github.com/spf13/cobra"
)

// TagsOptions - TODO
type TagsOptions struct{}

// Tags - TODO
type Tags struct {
	TagsOptions
}

// NewTags - TODO
func NewTags() *cobra.Command {
	var tags Tags

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "tags",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return tags.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run - TODO
func (t *Tags) Run(ctx context.Context, args []string) error {
	return nil
}
