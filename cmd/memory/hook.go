package memory

import (
	"context"
	"fmt"
	"io"

	"github.com/jamesl33/zk/internal/memory"
	"github.com/spf13/cobra"
)

// NewHook creates a new command for running the hooks.
func NewHook() *cobra.Command {
	cmd := cobra.Command{
		Short: "Runs a hook, reading its input from stdin",
		Use:   "hook",
	}

	cmd.AddCommand(
		NewHookExtract(),
		NewHookRecall(),
		NewHookStart(),
	)

	return &cmd
}

// hookCommand wraps fn as a hook which fails open: the agent carries on whatever goes wrong, and the error is logged in
// the vault.
func hookCommand(short, use string, fn func(context.Context, memory.Input, io.Writer) error) *cobra.Command {
	run := func(cmd *cobra.Command, _ []string) error {
		if memory.Saving() {
			return nil // The headless save run mustn't trigger hooks
		}

		err := runHook(cmd, fn)
		if err != nil {
			logHook(use, err)
		}

		return nil
	}

	cmd := cobra.Command{
		Short: short,
		Use:   use,
		RunE:  run,
	}

	return &cmd
}

// runHook reads the hook's input and runs fn with it.
func runHook(cmd *cobra.Command, fn func(context.Context, memory.Input, io.Writer) error) error {
	in, err := memory.ReadInput(cmd.InOrStdin())
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	return fn(cmd.Context(), in, cmd.OutOrStdout())
}

// logHook records a hook's failure in the vault, when there is one.
func logHook(use string, err error) {
	vault, verr := memory.VaultPath()
	if verr != nil {
		return
	}

	memory.Logf(vault, "%s: %v", use, err)
}
