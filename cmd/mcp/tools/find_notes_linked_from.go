package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
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

	var found []*note.Note

	err = notes.LinkedFrom(ctx, n, func(n *note.Note) {
		found = append(found, n)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find linked notes: %w", err)
	}

	output := FindNotesLinkedFromOutput{
		Notes: found,
	}

	return nil, &output, nil
}
