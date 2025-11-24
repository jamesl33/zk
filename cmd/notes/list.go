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

// ListOptions - TODO
//
// TODO (jamesl33): Add support for case-insensitive listing.
// TODO (jamesl33): Don't use boolean flags, use string flags then argument can be for the sub-directory.
type ListOptions struct {
	// Fixed - TODO
	Fixed string

	// Glob - TODO
	Glob string

	// Regex - TODO
	Regex string
}

// List - TODO
type List struct {
	ListOptions
}

// NewList - TODO
func NewList() *cobra.Command {
	var list List

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "list",
		// TODO
		Args: cobra.MaximumNArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return list.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringVar(
		&list.Fixed,
		"fixed",
		"",
		// TODO
		"",
	)

	cmd.Flags().StringVar(
		&list.Glob,
		"glob",
		"",
		// TODO
		"",
	)

	cmd.Flags().StringVar(
		&list.Regex,
		"regex",
		"",
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (l *List) Run(ctx context.Context, args []string) error {
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

	g.Go(func() error { defer w.Close(); return l.list(gctx, path, w) })

	_, err := io.Copy(os.Stdout, r)

	// TODO
	if errors.Is(err, syscall.EPIPE) {
		return nil
	}

	return err
}

// list - TODO
//
// TODO (jamesl33): Include the name in this search?
func (l *List) list(ctx context.Context, path string, w io.Writer) error {
	title, err := matcher.Title(l.Fixed, l.Glob, l.Regex)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(title),
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
