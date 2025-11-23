package notes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/gobwas/glob"
	icolor "github.com/jamesl33/zk/internal/color"
	"github.com/jamesl33/zk/internal/notes/searcher"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// SearchOptions - TODO
type SearchOptions struct {
	// Fixed - TODO
	Fixed bool

	// Glob - TODO
	Glob bool

	// Regex - TODO
	Regex bool
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
		Short: "",
		// TODO
		Use: "search",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return search.Run(cmd.Context(), args) },
	}

	cmd.Flags().BoolVar(
		&search.Fixed,
		"fixed",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&search.Glob,
		"glob",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&search.Regex,
		"regex",
		false,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (s *Search) Run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGPIPE)
	defer cancel()

	query := "."

	if len(args) >= 1 {
		query = args[0]
	}

	// TODO
	var (
		r, w    = io.Pipe()
		g, gctx = errgroup.WithContext(ctx)
	)

	g.Go(func() error { defer w.Close(); return s.search(gctx, query, w) })

	_, err := io.Copy(os.Stdout, r)

	// TODO
	if errors.Is(err, syscall.EPIPE) {
		return nil
	}

	return err
}

// search -  TODO
//
// TODO (jamesl33): Add human readable output.
func (s *Search) search(ctx context.Context, query string, w io.Writer) error {
	filter, err := s.filter(query)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	searcher, err := searcher.NewSearcher(
		searcher.WithPath("."),
		filter,
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	for n, err := range searcher.Many(ctx) {
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		fm, err := n.Frontmatter()
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		rel, err := filepath.Rel(".", n.Path)
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		fmt.Fprintf(w, "%s\x00(%s)\x00%s\n", icolor.Yellow(fm.Title), icolor.Blue(rel), n.Name())
	}

	return nil
}

// filter - TODO
func (s *Search) filter(query string) (func(*searcher.Options), error) {
	if s.Fixed {
		return searcher.WithFixed(query), nil
	}

	if s.Glob {
		return s.glob(query)
	}

	if s.Regex {
		return s.regex(query)
	}

	return searcher.WithFixed(query), nil
}

// glob - TODO
func (s *Search) glob(query string) (func(*searcher.Options), error) {
	parsed, err := glob.Compile("*" + query + "*")
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return searcher.WithGlob(parsed), nil
}

// regex - TODO
func (s *Search) regex(query string) (func(*searcher.Options), error) {
	parsed, err := regexp.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return searcher.WithRegex(parsed), nil
}
