package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vault"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ReadNoteInput defines the input for the ReadNote tool.
type ReadNoteInput struct {
	// Path is the path to the note to read.
	Path string `json:"path" jsonschema:"The path to a note"`
}

// ReadNoteOutput defines the output for the ReadNote tool.
type ReadNoteOutput struct {
	// Note is the note that was read, including its frontmatter.
	Note *note.Note `json:"note" jsonschema:"The note, including its frontmatter"`

	// Body is the note's content, without frontmatter.
	Body string `json:"body" jsonschema:"The note's content, without frontmatter"`
}

// ReadNote reads a note, returning its frontmatter and body.
func ReadNote(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *ReadNoteInput,
) (*mcp.CallToolResult, *ReadNoteOutput, error) {
	if _, err := vault.RootAbs("."); err != nil {
		return nil, nil, err
	}

	n, err := note.New(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open note: %w", err)
	}

	body, err := n.GetBody()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get body: %w", err)
	}

	output := ReadNoteOutput{
		Note: n,
		Body: body,
	}

	return nil, &output, nil
}
