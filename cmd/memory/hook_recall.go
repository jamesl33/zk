package memory

import (
	"context"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/jamesl33/zk/internal/memory"
	"github.com/spf13/cobra"
)

// recallTimeout bounds the embedding round trip a prompt waits on.
const recallTimeout = 3 * time.Second

// HookRecallOptions defines the options for the hook recall command.
type HookRecallOptions struct{}

// HookRecall defines the struct for the hook recall command.
type HookRecall struct {
	HookRecallOptions
}

// NewHookRecall creates a new command for the prompt hook, which injects the notes relevant to the prompt.
func NewHookRecall() *cobra.Command {
	var recall HookRecall

	return hookCommand("Injects the notes relevant to the prompt", "recall", recall.Run)
}

// Run injects the notes relevant to the prompt.
func (h *HookRecall) Run(ctx context.Context, in memory.Input, w io.Writer) error {
	vault, err := enter()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, recallTimeout)
	defer cancel()

	injected, err := memory.Recall(ctx, vault, in.SessionID, in.Prompt, memory.Find, maxDistance())
	if err != nil || injected == "" {
		return err
	}

	return write(w, memory.Context("UserPromptSubmit", injected))
}

// maxDistance is the distance within which a note is recalled.
func maxDistance() float32 {
	distance, err := strconv.ParseFloat(os.Getenv("ZK_MEMORY_MAX_DISTANCE"), 32)
	if err != nil || distance <= 0 {
		return memory.MaxDistance
	}

	return float32(distance)
}
