package tools

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/jamesl33/zk/internal/note"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// validNoteTypes are the note types that CreateNote will accept.
var validNoteTypes = []note.Type{"bibliographic", "fleeting", "index", "literature", "permanent"}

// CreateNoteInput defines the input for the CreateNote tool.
type CreateNoteInput struct {
	// Type is the note's type.
	Type string `json:"type" jsonschema:"The notes type, one of: bibliographic, fleeting, index, literature, permanent"`

	// Title is the note's title.
	Title string `json:"title" jsonschema:"The title for the note"`

	// Path is the directory the note should be created in.
	Path string `json:"path" jsonschema:"The directory to create the note in (e.g. '0 Inbox' for fleeting notes, '5 Bibliography' for bibliographic notes)"`

	// Tags are the note's tags.
	Tags []string `json:"tags" jsonschema:"The notes tags (e.g. short, simple, snake_case keywords to improve discoverability)"`
}

// CreateNoteOutput defines the output for the CreateNote tool.
type CreateNoteOutput struct {
	// Note is the note that was created.
	Note *note.Note `json:"note" jsonschema:"The note that was created"`
}

// CreateNote creates a new note.
func CreateNote(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *CreateNoteInput,
) (*mcp.CallToolResult, *CreateNoteOutput, error) {
	if !slices.Contains(validNoteTypes, note.Type(input.Type)) {
		return nil, nil, fmt.Errorf("invalid note type: %q", input.Type)
	}

	tags := input.Tags
	if tags == nil {
		tags = make([]string, 0)
	}

	n := note.Note{
		Path: note.Path(input.Path),
		Frontmatter: note.Frontmatter{
			Type:  note.Type(input.Type),
			Title: input.Title,
			Date:  time.Now().Format("2006-01-02"),
			Tags:  tags,
		},
	}

	err := n.Create()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to write note: %w", err)
	}

	output := CreateNoteOutput{
		Note: &n,
	}

	return nil, &output, nil
}
