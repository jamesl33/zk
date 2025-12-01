package note

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
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

// Run the command to find linked notes.
func (l *Links) Run(ctx context.Context, path string) error {
	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("failed to open note: %w", err)
	}

	var (
		// to if it's enabled by the user, or neither are enabled
		to = l.To || !(l.To || l.From)

		// from if it's enabled by the user, or neither are enabled
		from = l.From || !(l.To || l.From)
	)

	// Assign afterwards, as the values are part of both boolean expressions
	l.To = to
	l.From = from

	err = l.to(ctx, n)
	if err != nil {
		return fmt.Errorf("failed to list incoming notes: %w", err)
	}

	err = l.from(ctx, n)
	if err != nil {
		return fmt.Errorf("failed to list outgoing notes: %w", err)
	}

	return nil
}

// to lists the notes which link to the provided note.
func (l *Links) to(ctx context.Context, n *note.Note) error {
	// Not enabled, skip
	if !l.To {
		return nil
	}

	var (
		// name of the note, escaped for use in regular expressions
		name = regexp.QuoteMeta(n.Name())

		// pattern which matches links to this note
		pattern = fmt.Sprintf(`\[\[%s(\|.*?)?\]\]`, name)
	)

	matcher, err := matcher.Body("", "", pattern)
	if err != nil {
		return fmt.Errorf("failed to create matcher: %w", err)
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		fmt.Println(n.String0())
	}))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	return nil
}

// from lists the notes which link from the provided note.
func (l *Links) from(ctx context.Context, n *note.Note) error {
	// Not enabled, skip
	if !l.From {
		return nil
	}

	matchers := hs.Map(n.Links(), func(n string) matcher.Matcher { return matcher.Name(n) })

	// Must check for no matchers, as the default is to list all
	if len(matchers) == 0 {
		return nil
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		fmt.Println(n.String0())
	}))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	return nil
}
