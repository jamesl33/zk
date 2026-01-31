package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FindNotesLinkedFromInput defines the input for the FindNotesLinkedFrom tool.
type FindNotesLinkedFromInput struct {
	// Path is the path to the note to find linked notes for.
	Path string `json:"path" jsonschema:"The path to a note"`
}

// FindNotesLinkedFromOutput defines the output for the FindNotesLinkedFrom tool.
type FindNotesLinkedFromOutput struct {
	// Notes is a list of notes which are linked from the given note.
	Notes []*note.Note `json:"notes" jsonschema:"Notes which are linked from the given note"`
}

// FindNotesLinkedFrom finds notes which are linked from the given note.
func FindNotesLinkedFrom(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *FindNotesLinkedFromInput,
) (*mcp.CallToolResult, *FindNotesLinkedFromOutput, error) {
	n, err := note.New(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open note: %w", err)
	}

	notes, err := from(ctx, n)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find linked notes: %w", err)
	}

	output := FindNotesLinkedFromOutput{
		Notes: notes,
	}

	return nil, &output, nil
}

// from finds notes which are linked from the given note.
func from(ctx context.Context, n *note.Note) ([]*note.Note, error) {
	matchers := hs.Map(n.Links(), func(n string) matcher.Matcher { return matcher.Name(n) })

	// Must check for no matchers, as the default is to list all
	if len(matchers) == 0 {
		return nil, nil
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	found := make([]*note.Note, 0)

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		n.Body = ""
		found = append(found, n)
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	return found, nil
}
