package tools

import (
	"context"
	"fmt"
	"regexp"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FindNotesLinkedToInput defines the input for the FindNotesLinkedTo tool.
type FindNotesLinkedToInput struct {
	// Path is the path to the note to find linked notes for.
	Path string `json:"path" jsonschema:"The path to a note"`
}

// FindNotesLinkedToOutput defines the output for the FindNotesLinkedTo tool.
type FindNotesLinkedToOutput struct {
	// Notes is a list of notes which link to the given note.
	Notes []*note.Note `json:"notes" jsonschema:"Notes which link to the given note"`
}

// FindNotesLinkedTo finds notes which link to the given note.
func FindNotesLinkedTo(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *FindNotesLinkedToInput,
) (*mcp.CallToolResult, *FindNotesLinkedToOutput, error) {
	n, err := note.New(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open note: %w", err)
	}

	notes, err := to(ctx, n)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find linked notes: %w", err)
	}

	output := FindNotesLinkedToOutput{
		Notes: notes,
	}

	return nil, &output, nil
}

// to finds notes which link to the given note.
func to(ctx context.Context, n *note.Note) ([]*note.Note, error) {
	var (
		// name of the note, escaped for use in regular expressions
		name = regexp.QuoteMeta(n.Name())

		// pattern which matches links to this note
		pattern = fmt.Sprintf(`\[\[%s(\|.*?)?\]\]`, name)
	)

	matcher, err := matcher.Body("", "", pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create matcher: %w", err)
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher),
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
