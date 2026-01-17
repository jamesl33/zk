package note

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jamesl33/zk/internal/ai"
	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/regex"
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

	if len(n.Body) == 0 {
		return nil
	}

	client, err := ai.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	err = s.rewriteLinks(ctx, n)
	if err != nil {
		return fmt.Errorf("failed to rewrite links: %w", err)
	}

	prompt := `

%s

Without changing the meaning, produce a single sentence summary of the above note.`

	prompt = fmt.Sprintf(prompt, n.Body)

	// TODO (jamesl33): Handle the case where the model fails to summarize.
	content, err := client.Generate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate tags: %w", err)
	}

	fmt.Println(s.wrap(content))

	return nil
}

// rewriteLinks rewrites links within a note, converting the link into the title of the linked note; this is to enrich
// the context for summarization.
func (s *Summarize) rewriteLinks(ctx context.Context, n *note.Note) error {
	links := n.Links()

	if len(links) == 0 {
		return nil
	}

	matchers := hs.Map(links, func(n string) matcher.Matcher { return matcher.Name(n) })

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	titles := make(map[string]string)

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		titles[n.Name()] = n.Frontmatter.Title
	}))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	// rpl replaces links within notes, with the target notes title.
	rpl := func(match string) string {
		var (
			submatches = regex.Link.FindStringSubmatch(match)
			link       = submatches[regex.Link.SubexpIndex("link")]
		)

		if title, ok := titles[link]; ok {
			return title
		}

		return link
	}

	// Rewrite the note body
	n.Body = regex.Link.ReplaceAllStringFunc(n.Body, rpl)

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
