package memory

import (
	"context"
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
	return &cobra.Command{
		Short: short,
		Use:   use,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if memory.Saving() {
				return nil // The headless save run mustn't trigger hooks
			}

			in, err := memory.ReadInput(cmd.InOrStdin())
			if err != nil {
				return nil
			}

			err = fn(cmd.Context(), in, cmd.OutOrStdout())
			if err != nil {
				if vault, verr := memory.VaultPath(); verr == nil {
					memory.Logf(vault, "%s: %v", use, err)
				}
			}

			return nil
		},
	}
}
