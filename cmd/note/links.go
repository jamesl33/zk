package note

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
	"github.com/spf13/cobra"
)

// LinksOptions defines the options for the links command.
type LinksOptions struct {
	// To determines whether to show notes which link to the given note.
	To bool

	// From determines whether to show notes which are linked from the given note.
	From bool
}

// Links defines the struct for the links command.
type Links struct {
	LinksOptions
}

// NewLinks creates a new command for listing links to/from a note.
func NewLinks() *cobra.Command {
	var links Links

	cmd := cobra.Command{
		Short: "List the links to/from a note",
		Use:   "links <path>",
		Args:  cobra.ExactArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return links.Run(cmd.Context(), args[0]) },
	}

	cmd.Flags().BoolVar(
		&links.To,
		"to",
		false,
		"Display notes linked to the provied note",
	)

	cmd.Flags().BoolVar(
		&links.From,
		"from",
		false,
		"Display notes linked from the provied note",
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
	l.To, l.From = to, from

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

	err := notes.LinkedTo(ctx, n, func(n *note.Note) {
		fmt.Println(n.String())
	})
	if err != nil {
		return fmt.Errorf("failed to list incoming notes: %w", err)
	}

	return nil
}

// from lists the notes which link from the provided note.
func (l *Links) from(ctx context.Context, n *note.Note) error {
	// Not enabled, skip
	if !l.From {
		return nil
	}

	err := notes.LinkedFrom(ctx, n, func(n *note.Note) {
		fmt.Println(n.String())
	})
	if err != nil {
		return fmt.Errorf("failed to list outgoing notes: %w", err)
	}

	return nil
}
