package tools

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNote(t *testing.T) {
	tmp := withVault(t)

	input := CreateNoteInput{
		Type:  "permanent",
		Title: "My Note",
		Path:  "1 Projects",
		Tags:  []string{"tag1"},
		Body:  "Some content",
	}

	_, output, err := CreateNote(t.Context(), nil, &input)
	require.NoError(t, err)
	assert.Equal(t, "My Note", output.Note.Frontmatter.Title)

	matches, err := filepath.Glob(filepath.Join(tmp, "1 Projects", "*.md"))
	require.NoError(t, err)
	assert.Len(t, matches, 1)
}

func TestCreateNoteNilTags(t *testing.T) {
	withVault(t)

	input := CreateNoteInput{
		Type:  "fleeting",
		Title: "My Note",
		Path:  "0 Inbox",
	}

	_, output, err := CreateNote(t.Context(), nil, &input)
	require.NoError(t, err)
	assert.Empty(t, output.Note.Frontmatter.Tags)
	assert.NotNil(t, output.Note.Frontmatter.Tags)
}

func TestCreateNoteInvalidType(t *testing.T) {
	withVault(t)

	input := CreateNoteInput{
		Type:  "not-a-type",
		Title: "My Note",
		Path:  "0 Inbox",
	}

	_, _, err := CreateNote(t.Context(), nil, &input)
	assert.Error(t, err)
}

func TestCreateNoteNoVault(t *testing.T) {
	withEmptyDir(t)

	input := CreateNoteInput{
		Type:  "permanent",
		Title: "My Note",
		Path:  "1 Projects",
	}

	_, _, err := CreateNote(t.Context(), nil, &input)
	assert.Error(t, err)
}
