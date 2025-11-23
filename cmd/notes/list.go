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
	"github.com/jamesl33/zk/internal/notes/lister"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// ListOptions - TODO
//
// TODO (jamesl33): Add support for case-insensitive listing.
type ListOptions struct {
	// Fixed - TODO
	Fixed bool

	// Glob - TODO
	Glob bool

	// Regex - TODO
	Regex bool
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
		&list.Glob,
		"glob",
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

	return &cmd
}

// Run - TODO
func (l *List) Run(ctx context.Context, args []string) error {
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

	g.Go(func() error { defer w.Close(); return l.list(gctx, query, w) })

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
func (l *List) list(ctx context.Context, query string, w io.Writer) error {
	filter, err := l.filter(query)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		filter,
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	for n, err := range lister.Many(ctx) {
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
func (l *List) filter(query string) (func(*lister.Options), error) {
	if l.Fixed {
		return lister.WithFixed(query), nil
	}

	if l.Glob {
		return l.glob(query)
	}

	if l.Regex {
		return l.regex(query)
	}

	return lister.WithPath(query), nil
}

// glob - TODO
func (l *List) glob(query string) (func(*lister.Options), error) {
	parsed, err := glob.Compile("*" + query + "*")
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return lister.WithGlob(parsed), nil
}

// regex - TODO
func (l *List) regex(query string) (func(*lister.Options), error) {
	parsed, err := regexp.Compile(query)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return lister.WithRegex(parsed), nil
}
