package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FindRelatedNotesInput defines the input for the FindRelatedNotes tool.
type FindRelatedNotesInput struct {
	// Path is the path to the note to find related notes for.
	Path string `json:"path" jsonschema:"The path to a note"`
}

// FindRelatedNotesOutput defines the output for the FindRelatedNotes tool.
type FindRelatedNotesOutput struct {
	// Notes is a list of notes which are semantically similar to the given note.
	Notes []*note.Note `json:"notes" jsonschema:"Notes which are semantically similar to the given note"`
}

// FindRelatedNotes finds notes which are semantically similar to the given note.
func FindRelatedNotes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *FindRelatedNotesInput,
) (*mcp.CallToolResult, *FindRelatedNotesOutput, error) {
	n, err := note.New(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open note: %w", err)
	}

	var found []*note.Note

	err = notes.Find(ctx, n, hs.Infallible(func(n *note.Note) {
		found = append(found, n)
	}))
	if err != nil {
		return nil, nil, err
	}

	output := FindRelatedNotesOutput{
		Notes: found,
	}

	return nil, &output, nil
}
