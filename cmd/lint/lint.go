package lint

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/linter"
	"github.com/jamesl33/zk/internal/vault"
	"github.com/spf13/cobra"
)

// LintOptions defines the options for the lint command.
type LintOptions struct {
	// Archives includes the '4 Archives' directory in the linting results; it's excluded by
	// default.
	Archives bool
}

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

	cmd.Flags().BoolVar(&lint.Archives, "archives", false, "Include the '4 Archives' directory in the results")

	return &cmd
}

// Run lints the notes, printing warnings/errors.
func (l *Lint) Run(ctx context.Context, args []string) error {
	if _, err := vault.RootAbs("."); err != nil {
		return err
	}

	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	var opts []func(*linter.LintOptions)
	if l.Archives {
		opts = append(opts, linter.WithArchives())
	}

	errors, err := linter.NewLinter().Lint(ctx, path, opts...)
	if err != nil {
		return fmt.Errorf("failed to lint notes: %w", err)
	}

	for _, err := range errors {
		switch err.Line {
		case 0:
			fmt.Printf("%q: %s\n", err.Path, err.Message)
		default:
			fmt.Printf("%q:%d:%d: %s\n", err.Path, err.Line, err.Column, err.Message)
		}
	}

	return nil
}
