package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vector"
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
//
// TODO (jamesl33): De-duplicate this code?
func FindRelatedNotes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *FindRelatedNotesInput,
) (*mcp.CallToolResult, *FindRelatedNotesOutput, error) {
	n, err := note.New(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open note: %w", err)
	}

	db, err := vector.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	err = populate(ctx, db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to populate database: %w", err)
	}

	notes, err := db.Find(ctx, n)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find related notes: %w", err)
	}

	for _, n := range notes {
		n.Body = ""
	}

	output := FindRelatedNotesOutput{
		Notes: notes,
	}

	return nil, &output, nil
}
