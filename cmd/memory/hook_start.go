package memory

import (
	"context"
	"fmt"
	"io"

	"github.com/jamesl33/zk/internal/memory"
	"github.com/spf13/cobra"
)

// HookStartOptions defines the options for the hook start command.
type HookStartOptions struct{}

// HookStart defines the struct for the hook start command.
type HookStart struct {
	HookStartOptions
}

// NewHookStart creates a new command for the session start hook, which injects the digest.
func NewHookStart() *cobra.Command {
	var start HookStart

	return hookCommand("Injects the digest, and warns while built-in auto-memory is still on", "start", start.Run)
}

// Run injects the digest, and warns while built-in auto-memory is still on.
func (h *HookStart) Run(ctx context.Context, in memory.Input, w io.Writer) error {
	vault, err := enter()
	if err != nil {
		return err
	}

	if in.Source == "compact" {
		err = memory.ResetRecalled(vault, in.SessionID)
		if err != nil {
			return fmt.Errorf("failed to reset recalled notes: %w", err)
		}
	}

	_ = memory.Prune(vault)

	digest, err := memory.Digest(ctx, memory.ProjectKey(in.CWD))
	if err != nil {
		return err
	}

	out := memory.Context("SessionStart", digest)

	if memory.AutoMemoryEnabled(in.CWD) {
		out.SystemMessage = `zk-memory: built-in auto-memory is still on. Set "autoMemoryEnabled": false in ~/.claude/settings.json.`
	}

	return write(w, out)
}
