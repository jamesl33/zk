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

	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	// TODO
	var (
		r, w    = io.Pipe()
		g, gctx = errgroup.WithContext(ctx)
	)

	g.Go(func() error { defer w.Close(); return t.list(gctx, path, w) })

	_, err := io.Copy(os.Stdout, r)

	// TODO
	if errors.Is(err, syscall.EPIPE) {
		return nil
	}

	return err
}

// list - TODO
func (t *Tagged) list(ctx context.Context, path string, w io.Writer) error {
	tags, err := matcher.Tags(t.With, t.Without)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(tags),
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
