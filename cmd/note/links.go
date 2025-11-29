package note

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/spf13/cobra"
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
		RunE: func(cmd *cobra.Command, args []string) error { return links.Run(cmd.Context(), args[0]) },
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
func (l *Links) Run(ctx context.Context, path string) error {
	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	var (
		// TODO
		to = l.To || !(l.To || l.From)

		// TODO
		from = l.From || !(l.To || l.From)
	)

	// TODO
	l.To = to

	// TODO
	l.From = from

	err = l.to(ctx, n)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = l.from(ctx, n)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// to - TODO
func (l *Links) to(ctx context.Context, n *note.Note) error {
	// TODO
	if !l.To {
		return nil
	}

	var (
		// TODO
		name = regexp.QuoteMeta(n.Name())

		// TODO
		pattern = fmt.Sprintf(`\[\[%s(\|.*?)?\]\]`, name)
	)

	matcher, err := matcher.Body("", "", pattern)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		fmt.Println(n.String0())
	}))
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// from - TODO
func (l *Links) from(ctx context.Context, n *note.Note) error {
	// TODO
	if !l.From {
		return nil
	}

	matchers := hs.Map(n.Links(), func(n string) matcher.Matcher { return matcher.Name(n) })

	// TODO
	if len(matchers) == 0 {
		return nil
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		fmt.Println(n.String0())
	}))
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}
