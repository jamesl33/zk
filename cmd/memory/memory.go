// Package memory provides the hidden 'zk memory' commands, which Claude Code's hooks call to back Claude's persistent
// memory with a zk vault.
package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/jamesl33/zk/cmd/mcp"
	"github.com/jamesl33/zk/internal/memory"
	"github.com/spf13/cobra"
)

// recallTimeout bounds the embedding round trip a prompt waits on.
const recallTimeout = 3 * time.Second

// NewMemory creates the parent of the memory commands.
func NewMemory() *cobra.Command {
	cmd := cobra.Command{
		Short:  "Backs an agent's persistent memory with a vault",
		Use:    "memory",
		Hidden: true,
	}

	hook := cobra.Command{
		Short: "Runs a Claude Code hook, reading its input from stdin",
		Use:   "hook",
	}

	hook.AddCommand(
		hookCommand("start", Start),
		hookCommand("recall", Recall),
		hookCommand("extract", Detached),
	)

	cmd.AddCommand(&hook, NewSave(), NewMCP())

	return &cmd
}

// hookCommand wraps fn as a hook which fails open: Claude Code carries on whatever goes wrong, and the error is logged
// in the vault.
func hookCommand(name string, fn func(context.Context, memory.Input, io.Writer) error) *cobra.Command {
	return &cobra.Command{
		Short: "The " + name + " hook",
		Use:   name,
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
					memory.Logf(vault, "%s: %v", name, err)
				}
			}

			return nil
		},
	}
}

// enter ensures the memory vault and runs in it, returning its path.
func enter() (string, error) {
	vault, err := memory.VaultPath()
	if err != nil {
		return "", err
	}

	err = memory.EnsureVault(vault)
	if err != nil {
		return "", fmt.Errorf("failed to ensure vault: %w", err)
	}

	err = os.Chdir(vault)
	if err != nil {
		return "", fmt.Errorf("failed to enter vault: %w", err)
	}

	return vault, nil
}

func write(w io.Writer, out memory.Output) error {
	return json.NewEncoder(w).Encode(out)
}

// Start injects the digest, and warns while built-in auto-memory is still on.
func Start(ctx context.Context, in memory.Input, w io.Writer) error {
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

// Recall injects the notes relevant to the prompt.
func Recall(ctx context.Context, in memory.Input, w io.Writer) error {
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

// Detached hands the transcript to a detached save, so compaction and exit don't wait.
func Detached(_ context.Context, in memory.Input, _ io.Writer) error {
	if in.TranscriptPath == "" || in.SessionID == "" {
		return nil
	}

	return memory.Detach("memory", "save", "--session", in.SessionID, "--transcript", in.TranscriptPath, "--cwd", in.CWD)
}

// SaveOptions defines the options for the save command.
type SaveOptions struct {
	Session    string
	Transcript string
	CWD        string
}

// NewSave creates the command which saves a transcript's learnings.
func NewSave() *cobra.Command {
	var o SaveOptions

	cmd := cobra.Command{
		Short: "Saves the learnings of a session transcript using a headless Claude",
		Use:   "save",
		RunE:  func(cmd *cobra.Command, _ []string) error { return o.Run(cmd.Context()) },
	}

	cmd.Flags().StringVar(&o.Session, "session", "", "Session id")
	cmd.Flags().StringVar(&o.Transcript, "transcript", "", "Path to the transcript")
	cmd.Flags().StringVar(&o.CWD, "cwd", "", "Directory the session ran in")

	return &cmd
}

// Run saves the learnings, logging rather than returning errors as nothing is waiting for them.
func (o *SaveOptions) Run(ctx context.Context) error {
	vault, err := enter()
	if err != nil {
		return err
	}

	model := os.Getenv("ZK_MEMORY_MODEL")
	if model == "" {
		model = "haiku"
	}

	err = memory.Extract(ctx, memory.Extraction{
		Vault:      vault,
		Session:    o.Session,
		Transcript: o.Transcript,
		Key:        memory.ProjectKey(o.CWD),
		Run:        memory.ClaudeRunner(model),
	})
	if err != nil {
		memory.Logf(vault, "save: %v", err)
	}

	return nil
}

// NewMCP creates the command which serves the memory vault over the model context protocol. Entering the vault here
// (rather than passing '--vault') creates it first, as the server may start before the SessionStart hook has run.
func NewMCP() *cobra.Command {
	return &cobra.Command{
		Short: "Runs the model context protocol server against the memory vault",
		Use:   "mcp",
		RunE: func(cmd *cobra.Command, _ []string) error {
			_, err := enter()
			if err != nil {
				return err
			}

			var server mcp.MCP

			return server.Run(cmd.Context())
		},
	}
}

// maxDistance is the distance within which a note is recalled.
func maxDistance() float32 {
	distance, err := strconv.ParseFloat(os.Getenv("ZK_MEMORY_MAX_DISTANCE"), 32)
	if err != nil || distance <= 0 {
		return memory.MaxDistance
	}

	return float32(distance)
}
