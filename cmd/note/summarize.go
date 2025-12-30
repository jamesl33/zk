package note

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jamesl33/zk/internal/ai"
	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// SummarizeOptions - TODO
type SummarizeOptions struct{}

// Summarize - TODO
type Summarize struct {
	SummarizeOptions
}

// NewSummarize - TODO
func NewSummarize() *cobra.Command {
	var summarize Summarize

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "summarize",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return summarize.Run(cmd.Context(), args[0]) },
	}

	return &cmd
}

// Run the command to find linked notes.
func (s *Summarize) Run(ctx context.Context, path string) error {
	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("failed to open note: %w", err)
	}

	// TODO
	if len(n.Body) == 0 {
		return nil
	}

	client, err := ai.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	prompt := `

%s

Without changing the meaning, produce a single sentence - less that 80 characters - summary of the above note.`

	prompt = fmt.Sprintf(prompt, n.Body)

	content, err := client.Generate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate tags: %w", err)
	}

	fmt.Println(content)

	return nil
}
