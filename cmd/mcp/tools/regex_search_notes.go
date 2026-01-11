package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegexSearchNotesInput - TODO
type RegexSearchNotesInput struct {
	// Path - TODO
	Path string `json:"path" jsonschema:"The path to the directory to search in"`

	// Expression - TODO
	Expression string `json:"expression" jsonschema:"The (RE2) regular expression to use"`
}

// RegexSearchNotesOutput - TODO
type RegexSearchNotesOutput struct {
	// Notes - TODO
	Notes []*note.Note `json:"notes" jsonschema:"Notes where the title, tags or content match the regular expression"`
}

// RegexSearchNotes - TODO
//
// TODO (jamesl33): De-duplicate this code?
func RegexSearchNotes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *RegexSearchNotesInput,
) (*mcp.CallToolResult, *RegexSearchNotesOutput, error) {
	pm, err := matcher.Path("", "", input.Expression)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create path matcher: %w", err)
	}

	entire, err := matcher.Entire("", "", input.Expression)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create entire matcher: %w", err)
	}

	lister, err := lister.NewLister(
		lister.WithPath(input.Path),
		lister.WithMatcher(matcher.Or(pm, entire)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create lister: %w", err)
	}

	found := make([]*note.Note, 0)

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		// TODO
		n.Body = ""

		// TODO
		found = append(found, n)
	}))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search notes: %w", err)
	}

	output := RegexSearchNotesOutput{
		Notes: found,
	}

	return nil, &output, nil
}
