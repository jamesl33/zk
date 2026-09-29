package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jamesl33/zk/internal/note"
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
	if err := checkPath(input.Path); err != nil {
		return nil, nil, err
	}

	// Only notes can be deleted, not other files which happen to be in the vault.
	if filepath.Ext(input.Path) != ".md" {
		return nil, nil, fmt.Errorf("path %q is not a markdown note", input.Path)
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
