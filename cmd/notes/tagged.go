package notes

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	icolor "github.com/jamesl33/zk/internal/color"
	"github.com/jamesl33/zk/internal/notes/lister"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// TaggedOptions - TODO
//
// TODO (jamesl33): Move this sub-command to `zk notes list tagged`?
type TaggedOptions struct {
	// With - TODO
	With []string

	// Without - TODO
	Without []string
}

// Tagged - TODO
type Tagged struct {
	TaggedOptions
}

// NewTagged - TODO
func NewTagged() *cobra.Command {
	var tagged Tagged

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "tagged",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return tagged.Run(cmd.Context(), args) },
	}

	cmd.Flags().StringArrayVar(
		&tagged.With,
		"with",
		nil,
		"",
	)

	cmd.Flags().StringArrayVar(
		&tagged.Without,
		"without",
		nil,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (t *Tagged) Run(ctx context.Context, args []string) error {
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

	g.Go(func() error { defer w.Close(); return t.list(gctx, query, w) })

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
func (t *Tagged) list(ctx context.Context, query string, w io.Writer) error {
	lister, err := lister.NewLister(
		lister.WithPath("."),
		t.filter(query),
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
func (t *Tagged) filter(query string) func(*lister.Options) {
	if len(t.With) != 0 {
		return lister.WithTagged(t.With)
	}

	return lister.WithNotTagged(t.Without)
}
