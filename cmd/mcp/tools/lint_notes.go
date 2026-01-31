package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/regex"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// LintError defines a single linting error.
type LintError struct {
	Path    string `json:"path" jsonschema:"The path to the note"`
	Message string `json:"message" jsonschema:"A description of the linting failure"`
}

// LintNotesInput defines the input for the LintNotes tool.
type LintNotesInput struct{}

// LintNotesOutput defines the output for the LintNotes tool.
type LintNotesOutput struct {
	// Errors is a list of linting errors that were found.
	Errors []*LintError `json:"errors" jsonschema:"A list of linting errors that were found."`
}

// LintNotes lints notes in a given directory.
//
// TODO (jamesl33): De-duplicate this logic.
func LintNotes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *LintNotesInput,
) (*mcp.CallToolResult, *LintNotesOutput, error) {
	ids := make([]string, 0)

	lstr, err := lister.NewLister(
		lister.WithPath("."),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lstr.Many(ctx), hs.Infallible(func(n *note.Note) {
		ids = append(ids, n.Name())
	}))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list notes: %w", err)
	}

	entire, err := matcher.Entire("", "", regex.Link.String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create entire matcher: %w", err)
	}

	lstr, err = lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(entire),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create lister: %w", err)
	}

	errors := make([]*LintError, 0)

	err = iterator.ForEach2(lstr.Many(ctx), hs.Infallible(func(n *note.Note) {
		for _, link := range hs.Difference(n.Links(), ids) {
			err := LintError{
				Path:    n.Path,
				Message: fmt.Sprintf("Link %q is broken (linkcheck)", link),
			}

			errors = append(errors, &err)
		}
	}))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list notes: %w", err)
	}

	output := LintNotesOutput{
		Errors: errors,
	}

	return nil, &output, nil
}
