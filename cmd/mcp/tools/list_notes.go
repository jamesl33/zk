package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/note"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListNotesInput defines the input for the ListNotes tool.
type ListNotesInput struct {
	// Path is a directory in which to list notes.
	Path string `json:"path" jsonschema:"A directory in which to list notes, use '.' to list everything"`
}

// ListNotesOutput defines the output for the ListNotes tool.
type ListNotesOutput struct {
	// Notes is a list of notes that were in the given directory.
	Notes []*note.Note `json:"notes" jsonschema:"Notes that were in the given directory"`
}

// ListNotes lists notes in a given directory.
func ListNotes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *ListNotesInput,
) (*mcp.CallToolResult, *ListNotesOutput, error) {
	lister, err := lister.NewLister(
		lister.WithPath(input.Path),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create lister: %w", err)
	}

	found := make([]*note.Note, 0)

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		n.Body = ""
		found = append(found, n)
	}))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search notes: %w", err)
	}

	output := ListNotesOutput{
		Notes: found,
	}

	return nil, &output, nil
}
