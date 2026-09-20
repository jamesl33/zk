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
	Line    int    `json:"line" jsonschema:"The 1-based line number the error occurs on, or 0 if it applies to the whole note"`
	Column  int    `json:"column" jsonschema:"The 1-based column the error occurs on, or 0 if line is also 0"`
}

// LintNotesInput defines the input for the LintNotes tool.
type LintNotesInput struct {
	// Path is the directory to lint.
	Path string `json:"path" jsonschema:"The directory to lint, use '.' to lint everything"`

	// Archives includes the '4 Archives' directory in the linting results; it's excluded by
	// default.
	Archives bool `json:"archives,omitempty" jsonschema:"Include the '4 Archives' directory in the results, false by default"`
}

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
	path := input.Path
	if path == "" {
		path = "."
	}

	var opts []func(*linter.LintOptions)
	if input.Archives {
		opts = append(opts, linter.WithArchives())
	}

	errors, err := linter.NewLinter().Lint(ctx, path, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to lint notes: %w", err)
	}

	output := LintNotesOutput{
		Errors: hs.Map(errors, func(err *linter.LintError) *LintError { return (*LintError)(err) }),
	}

	return nil, &output, nil
}
