package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/note"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ReadNoteInput - TODO
type ReadNoteInput struct {
	// Path - TODO
	Path string `json:"path" jsonschema:"The path to a note"`
}

// ReadNoteOutput - TODO
type ReadNoteOutput struct {
	// Note - TODO
	Note *note.Note `json:"notes" jsonschema:"The contents of the note, read from disk"`
}

// ReadNote - TODO
//
// TODO (jamesl33): De-duplicate this code?
func ReadNote(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *ReadNoteInput,
) (*mcp.CallToolResult, *ReadNoteOutput, error) {
	n, err := note.New(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open note: %w", err)
	}

	output := ReadNoteOutput{
		Note: n,
	}

	return nil, &output, nil
}
