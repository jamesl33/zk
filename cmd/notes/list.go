package notes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/jamesl33/zk/internal/notes"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// ListOptions - TODO
type ListOptions struct {
	// Fixed - TODO
	Fixed bool

	// Regex - TODO
	Regex bool

	// Glob - TODO
	Glob bool
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

	cmd.Flags().BoolVar(
		&list.Fixed,
		"fixed",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&list.Regex,
		"regex",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&list.Glob,
		"glob",
		false,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (l *List) Run(ctx context.Context, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGPIPE)
	defer cancel()

	var path string

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
// TODO (jamesl33): Add human readable output.
// TODO (jamesl33): Add support for listing with 'fd'.
func (l *List) list(ctx context.Context, path string, w io.Writer) error {
	for n, err := range notes.List(ctx, path) {
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		fm, err := n.Frontmatter()
		if err != nil {
			return fmt.Errorf("%w", err) // TODO
		}

		fmt.Fprintf(w, "%s\x00%s\n", fm.Title, n.Path)
	}

	return nil
}
