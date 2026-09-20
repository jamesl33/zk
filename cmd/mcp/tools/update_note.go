package tools

import (
	"context"
	"fmt"
	"slices"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vault"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// UpdateNoteInput defines the input for the UpdateNote tool.
type UpdateNoteInput struct {
	// Path is the path to the note to update.
	Path string `json:"path" jsonschema:"The path to the note to update"`

	// Type, if set, replaces the note's type.
	Type *string `json:"type,omitempty" jsonschema:"If set, replaces the notes type, one of: bibliographic, fleeting, index, literature, permanent"`

	// Title, if set, replaces the note's title.
	Title *string `json:"title,omitempty" jsonschema:"If set, replaces the notes title"`

	// Tags, if set, replaces the note's tags.
	Tags *[]string `json:"tags,omitempty" jsonschema:"If set, replaces the notes tags"`

	// Body, if set, replaces the note's body.
	Body *string `json:"body,omitempty" jsonschema:"If set, replaces the notes body"`
}

// UpdateNoteOutput defines the output for the UpdateNote tool.
type UpdateNoteOutput struct {
	// Note is the note as it was left after the update.
	Note *note.Note `json:"note" jsonschema:"The note, after being updated"`
}

// UpdateNote updates a note in place. Only the fields which are set on the input are changed;
// anything left unset is preserved as-is.
func UpdateNote(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *UpdateNoteInput,
) (*mcp.CallToolResult, *UpdateNoteOutput, error) {
	if _, err := vault.RootAbs("."); err != nil {
		return nil, nil, err
	}

	n, err := note.New(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open note: %w", err)
	}

	if input.Type != nil {
		if !slices.Contains(validNoteTypes, note.Type(*input.Type)) {
			return nil, nil, fmt.Errorf("invalid note type: %q", *input.Type)
		}

		n.Frontmatter.Type = note.Type(*input.Type)
	}

	if input.Title != nil {
		n.Frontmatter.Title = *input.Title
	}

	if input.Tags != nil {
		n.Frontmatter.Tags = *input.Tags
	}

	if input.Body != nil {
		n.SetBody(*input.Body)
	}

	err = n.Write()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write note: %w", err)
	}

	output := UpdateNoteOutput{
		Note: n,
	}

	return nil, &output, nil
}
