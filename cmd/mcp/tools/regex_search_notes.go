package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegexSearchNotesInput defines the input for the RegexSearchNotes tool.
type RegexSearchNotesInput struct {
	// Path is the path to the directory to search in.
	Path string `json:"path" jsonschema:"The path to the directory to search in, use '.' to search everything"`

	// Expression is the (RE2) regular expression to use.
	Expression string `json:"expression" jsonschema:"The (RE2) regular expression to use"`
}

// RegexSearchNotesOutput defines the output for the RegexSearchNotes tool.
type RegexSearchNotesOutput struct {
	// Notes is a list of notes where the title, tags or content match the regular expression.
	Notes []*note.Note `json:"notes" jsonschema:"Notes where the title, tags or content match the regular expression"`
}

// RegexSearchNotes finds notes where the title, tags or content match the regular expression.
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

	var found []*note.Note

	err = notes.Search(ctx, input.Path, matcher.Or(pm, entire), func(n *note.Note) {
		n.Body = ""
		found = append(found, n)
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to search notes: %w", err)
	}

	output := RegexSearchNotesOutput{
		Notes: found,
	}

	return nil, &output, nil
}
