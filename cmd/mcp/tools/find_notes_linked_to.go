package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
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

	var found []*note.Note

	err = notes.LinkedTo(ctx, n, func(n *note.Note) {
		found = append(found, n)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find linked notes: %w", err)
	}

	output := FindNotesLinkedToOutput{
		Notes: found,
	}

	return nil, &output, nil
}
