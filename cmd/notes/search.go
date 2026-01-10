package notes

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// SearchOptions - TODO
//
// TODO (jamesl33): Add support for case-insensitive search.
type SearchOptions struct {
	// Fixed - TODO
	Fixed string

	// Glob - TODO
	Glob string

	// Regex - TODO
	Regex string
}

// Search - TODO
type Search struct {
	SearchOptions
}

// NewSearch - TODO
func NewSearch() *cobra.Command {
	var search Search

	cmd := cobra.Command{
		// TODO
		Short: "Search the content of notes, listing the matching notes",
		// TODO
		Use: "search [directory]",
		// TODO
		Args: cobra.MaximumNArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return search.Run(cmd.Context(), args) },
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
//
// TODO (jamesl33): Include the tags in this search?
func (s *Search) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	title, err := matcher.Title(s.Fixed, s.Glob, s.Regex)
	if err != nil {
		return fmt.Errorf("failed to create title matcher: %w", err)
	}

	body, err := matcher.Body(s.Fixed, s.Glob, s.Regex)
	if err != nil {
		return fmt.Errorf("failed to create body matcher: %w", err)
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(matcher.Or(title, body)),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		fmt.Println(n.String0())
	}))
	if err != nil {
		return fmt.Errorf("failed to search notes: %w", err)
	}

	return nil
}
