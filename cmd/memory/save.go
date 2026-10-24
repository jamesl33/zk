package memory

import (
	"context"
	"os"

	"github.com/jamesl33/zk/internal/memory"
	"github.com/spf13/cobra"
)

// SaveOptions defines the options for the save command.
type SaveOptions struct {
	// Session is the id of the session.
	Session string

	// Transcript is the path to the session transcript.
	Transcript string

	// CWD is the directory the session ran in.
	CWD string
}

// Save defines the struct for the save command.
type Save struct {
	SaveOptions
}

// NewSave creates a new command for saving the learnings of a session transcript.
func NewSave() *cobra.Command {
	var save Save

	cmd := cobra.Command{
		Short: "Saves the learnings of a session transcript using a headless agent",
		Use:   "save",
		RunE:  func(cmd *cobra.Command, _ []string) error { return save.Run(cmd.Context()) },
	}

	cmd.Flags().StringVar(&save.Session, "session", "", "Session id")
	cmd.Flags().StringVar(&save.Transcript, "transcript", "", "Path to the transcript")
	cmd.Flags().StringVar(&save.CWD, "cwd", "", "Directory the session ran in")

	return &cmd
}

// Run saves the learnings, logging rather than returning errors as nothing is waiting for them.
func (s *Save) Run(ctx context.Context) error {
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
		Session:    s.Session,
		Transcript: s.Transcript,
		Key:        memory.ProjectKey(s.CWD),
		Run:        memory.ClaudeRunner(model),
	})
	if err != nil {
		memory.Logf(vault, "save: %v", err)
	}

	return nil
}
