package tools

import (
	"context"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SemanticSearchNotesInput defines the input for the SemanticSearchNotes tool.
type SemanticSearchNotesInput struct {
	// Query is a generic query, phrase, word or pattern used to find semantically similar notes.
	Query string `json:"query" jsonschema:"A generic query, phrase, word or pattern used to find semantically similar notes"`
}

// SemanticSearchNotesOutput defines the output for the SemanticSearchNotes tool.
type SemanticSearchNotesOutput struct {
	// Notes is a list of notes that are semantically similar to the query.
	Notes []*note.Note `json:"notes" jsonschema:"Notes that are semantically similar to the query"`
}

// SemanticSearchNotes finds notes that are semantically similar to the query.
func SemanticSearchNotes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *SemanticSearchNotesInput,
) (*mcp.CallToolResult, *SemanticSearchNotesOutput, error) {
	n := &note.Note{
		Body: input.Query,
	}

	var found []*note.Note

	err := notes.Find(ctx, n, hs.Infallible(func(n *note.Note) {
		n.Body = ""
		found = append(found, n)
	}))
	if err != nil {
		return nil, nil, err
	}

	output := SemanticSearchNotesOutput{
		Notes: found,
	}

	return nil, &output, nil
}
