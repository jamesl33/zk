package notes

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
	"github.com/spf13/cobra"
)

// SearchOptions defines the options for the search command.
//
// TODO (jamesl33): Add support for case-insensitive search.
type SearchOptions struct {
	// Fixed filters notes by title/content using a case-sensitive fixed-string search.
	Fixed string

	// Glob filters notes by title/content using a case-sensitive glob pattern.
	Glob string

	// Regex filters notes by title/content using a regular expression (RE2).
	Regex string
}

// Search defines the struct for the search command.
type Search struct {
	SearchOptions
}

// NewSearch creates a new command for searching notes.
func NewSearch() *cobra.Command {
	var search Search

	cmd := cobra.Command{
		Short: "Search the content of notes, listing the matching notes",
		Use:   "search [directory]",
		Args:  cobra.MaximumNArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return search.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringVar(
		&search.Fixed,
		"fixed",
		"",
		"Filter notes by title/content using a case-sensitive fixed-string search",
	)

	cmd.Flags().StringVar(
		&search.Glob,
		"glob",
		"",
		"Filter notes by title/content using a case-sensitive glob pattern",
	)

	cmd.Flags().StringVar(
		&search.Regex,
		"regex",
		"",
		"Filter notes by title/content using a regular expression (RE2)",
	)

	return &cmd
}

// Run searches notes (e.g. titles, bodies).
func (s *Search) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	pm, err := matcher.Path(s.Fixed, s.Glob, s.Regex)
	if err != nil {
		return fmt.Errorf("failed to create path matcher: %w", err)
	}

	entire, err := matcher.Entire(s.Fixed, s.Glob, s.Regex)
	if err != nil {
		return fmt.Errorf("failed to create entire matcher: %w", err)
	}

	err = notes.Search(ctx, path, matcher.Or(pm, entire), func(n *note.Note) {
		fmt.Println(n.String0())
	})
	if err != nil {
		return fmt.Errorf("failed to search notes: %w", err)
	}

	return nil
}
