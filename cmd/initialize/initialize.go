package initialize

import (
	"context"
	"fmt"
	"os"

	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed instructions.md
var instructions []byte

// InitializeOptions - TODO
type InitializeOptions struct{}

// Initialize - TODO
type Initialize struct {
	InitializeOptions
}

// NewInitialize - TODO
func NewInitialize() *cobra.Command {
	var index Initialize

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "initialize",
		// TODO
		RunE: func(cmd *cobra.Command, _ []string) error { return index.Run(cmd.Context()) },
	}

	return &cmd
}

// Run initialization.
func (i *Initialize) Run(ctx context.Context) error {
	err := os.WriteFile("GEMINI.md", instructions, 0o644)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}
