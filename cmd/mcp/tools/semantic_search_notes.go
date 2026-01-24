package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vector"
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
//
// TODO (jamesl33): De-duplicate this code?
func SemanticSearchNotes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *SemanticSearchNotesInput,
) (*mcp.CallToolResult, *SemanticSearchNotesOutput, error) {
	db, err := vector.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	err = populate(ctx, db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to populate database: %w", err)
	}

	n := &note.Note{
		Body: input.Query,
	}

	notes, err := db.Find(ctx, n)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find related notes: %w", err)
	}

	for _, n := range notes {
		n.Body = ""
	}

	output := SemanticSearchNotesOutput{
		Notes: notes,
	}

	return nil, &output, nil
}

// populate the index by updating embeddings for notes that have been updated.
func populate(ctx context.Context, db *vector.DB) error {
	lister, err := lister.NewLister(
		lister.WithPath("."),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), func(n *note.Note) error {
		return db.Upsert(ctx, n)
	})
	if err != nil {
		return fmt.Errorf("failed to upsert embeddings: %w", err)
	}

	return nil
}
