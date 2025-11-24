package notes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/jamesl33/zk/internal/notes/lister"
	"github.com/jamesl33/zk/internal/notes/matcher"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
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
		Short: "",
		// TODO
		Use: "search",
		// TODO
		Args: cobra.MaximumNArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return search.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringVar(
		&search.Fixed,
		"fixed",
		"",
		// TODO
		"",
	)

	cmd.Flags().StringVar(
		&search.Glob,
		"glob",
		"",
		// TODO
		"",
	)

	cmd.Flags().StringVar(
		&search.Regex,
		"regex",
		"",
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (s *Search) Run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGPIPE)
	defer cancel()

	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	// TODO
	var (
		r, w    = io.Pipe()
		g, gctx = errgroup.WithContext(ctx)
	)

	g.Go(func() error { defer w.Close(); return s.search(gctx, path, w) })

	_, err := io.Copy(os.Stdout, r)

	// TODO
	if errors.Is(err, syscall.EPIPE) {
		return nil
	}

	return err
}

// search -  TODO
func (s *Search) search(ctx context.Context, path string, w io.Writer) error {
	matcher, err := matcher.NewBody(s.Fixed, s.Glob, s.Regex)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(matcher),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	for n, err := range lister.Many(ctx) {
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		fmt.Fprintln(w, n.String0())
	}

	return nil
}
