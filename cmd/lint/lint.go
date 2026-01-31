package lint

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/linter"
	"github.com/spf13/cobra"
)

// LintOptions defines the options for the lint command.
type LintOptions struct{}

// Lint defines the struct for the lint command.
type Lint struct {
	LintOptions
}

// NewLint creates a new command for linting notes.
func NewLint() *cobra.Command {
	var lint Lint

	cmd := cobra.Command{
		Short: "Lints notes",
		Use:   "lint [directory]",
		Args:  cobra.MaximumNArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return lint.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run lints the notes, printing warnings/errors.
func (l *Lint) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	errors, err := linter.NewLinter().Lint(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to lint notes: %w", err)
	}

	for _, err := range errors {
		fmt.Printf("%q: %s\n", err.Path, err.Message)
	}

	return nil
}
