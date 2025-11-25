package note

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes/lister"
	"github.com/jamesl33/zk/internal/notes/matcher"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

// LinksOptions - TODO
type LinksOptions struct {
	// To - TODO
	To bool

	// From - TODO
	From bool
}

// Links - TODO
type Links struct {
	LinksOptions
}

// NewLinks - TODO
func NewLinks() *cobra.Command {
	var links Links

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "links",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return links.Run(cmd.Context(), args) },
	}

	cmd.Flags().BoolVar(
		&links.To,
		"to",
		false,
		// TODO
		"",
	)

	cmd.Flags().BoolVar(
		&links.From,
		"from",
		false,
		// TODO
		"",
	)

	return &cmd
}

// Run - TODO
func (l *Links) Run(ctx context.Context, args []string) error {
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

	g.Go(func() error { defer w.Close(); return l.links(gctx, path, w) })

	_, err := io.Copy(os.Stdout, r)

	// TODO
	if errors.Is(err, syscall.EPIPE) {
		return nil
	}

	return err
}

func (l *Links) links(ctx context.Context, path string, w io.Writer) error {
	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	matchers := hs.Map(n.Links(), func(n string) matcher.Matcher { return matcher.Name(n) })

	// TODO (jamesl33): This doesn't work; the path needs to be the root for the Zettelkasten. Viper?
	lister, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(matcher.Or(matchers...)),
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
