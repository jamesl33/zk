package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vault"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// DeleteNoteInput defines the input for the DeleteNote tool.
type DeleteNoteInput struct {
	// Path is the path to the note to delete.
	Path string `json:"path" jsonschema:"The path to the note to delete"`
}

// DeleteNoteOutput defines the output for the DeleteNote tool.
type DeleteNoteOutput struct{}

// DeleteNote deletes a note.
func DeleteNote(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input *DeleteNoteInput,
) (*mcp.CallToolResult, *DeleteNoteOutput, error) {
	if _, err := vault.RootAbs("."); err != nil {
		return nil, nil, err
	}

	n, err := note.New(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open note: %w", err)
	}

	err = os.Remove(n.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to delete note: %w", err)
	}

	return nil, &DeleteNoteOutput{}, nil
}
