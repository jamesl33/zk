package tools

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/linter"
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
func LintNotes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *LintNotesInput,
) (*mcp.CallToolResult, *LintNotesOutput, error) {
	errors, err := linter.NewLinter().Lint(ctx, ".")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lint notes: %w", err)
	}

	output := LintNotesOutput{
		Errors: hs.Map(errors, func(err *linter.LintError) *LintError { return (*LintError)(err) }),
	}

	return nil, &output, nil
}
