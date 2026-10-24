package memory

import (
	"context"
	"io"

	"github.com/jamesl33/zk/internal/memory"
	"github.com/spf13/cobra"
)

// HookExtractOptions defines the options for the hook extract command.
type HookExtractOptions struct{}

// HookExtract defines the struct for the hook extract command.
type HookExtract struct {
	HookExtractOptions
}

// NewHookExtract creates a new command for the compaction and session end hook, which saves what the session learned.
func NewHookExtract() *cobra.Command {
	var extract HookExtract

	return hookCommand("Hands the transcript to a detached save", "extract", extract.Run)
}

// Run hands the transcript to a detached save, so compaction and exit don't wait.
func (h *HookExtract) Run(_ context.Context, in memory.Input, _ io.Writer) error {
	if in.TranscriptPath == "" || in.SessionID == "" {
		return nil
	}

	return memory.Detach("memory", "save", "--session", in.SessionID, "--transcript", in.TranscriptPath, "--cwd", in.CWD)
}
