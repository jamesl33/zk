package note

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jamesl33/zk/internal/ai"
	"github.com/jamesl33/zk/internal/links"
	"github.com/jamesl33/zk/internal/note"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

// SummarizeOptions defines the options for the summarize command.
type SummarizeOptions struct{}

// Summarize defines the struct for the summarize command.
type Summarize struct {
	SummarizeOptions
}

// NewSummarize creates a new command for summarizing notes.
func NewSummarize() *cobra.Command {
	var summarize Summarize

	cmd := cobra.Command{
		Short: "Summarize a note",
		Use:   "summarize <path>",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return summarize.Run(cmd.Context(), args[0]) },
	}

	return &cmd
}

func (s *Summarize) Run(ctx context.Context, path string) error {
	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("failed to open note: %w", err)
	}

	err = links.Replace(ctx, n)
	if err != nil {
		return fmt.Errorf("failed to replace links: %w", err)
	}

	body, err := n.GetBody()
	if err != nil {
		return fmt.Errorf("failed to get body: %w", err)
	}

	if len(body) == 0 {
		return nil
	}

	client, err := ai.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	prompt := `

%s

Without changing the meaning, produce a single sentence summary of the above note.`

	prompt = fmt.Sprintf(prompt, body)

	// TODO (jamesl33): Handle the case where the model fails to summarize.
	content, err := client.Generate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate tags: %w", err)
	}

	fmt.Println(s.wrap(content))

	return nil
}

func (s *Summarize) wrap(content string) string {
	raw := os.Getenv("FZF_PREVIEW_COLUMNS")

	if raw == "" {
		return content
	}

	columns, _ := strconv.ParseUint(raw, 10, 64)

	if columns == 0 {
		return content
	}

	return wordwrap.WrapString(content, uint(columns))
}
