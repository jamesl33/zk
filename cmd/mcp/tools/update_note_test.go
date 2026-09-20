package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateNote(t *testing.T) {
	tmp := withVault(t)

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\ntype: fleeting\n---\nOld body"), 0o644)
	require.NoError(t, err)

	input := UpdateNoteInput{
		Path:  path,
		Type:  ptr.To("permanent"),
		Title: ptr.To("New Title"),
		Tags:  ptr.To([]string{"tag1"}),
		Body:  ptr.To("New body"),
	}

	_, output, err := UpdateNote(t.Context(), nil, &input)
	require.NoError(t, err)
	assert.Equal(t, note.Type("permanent"), output.Note.Frontmatter.Type)
	assert.Equal(t, "New Title", output.Note.Frontmatter.Title)
	assert.Equal(t, []string{"tag1"}, output.Note.Frontmatter.Tags)

	reloaded, err := note.New(path)
	require.NoError(t, err)
	assert.Equal(t, "New Title", reloaded.Frontmatter.Title)

	body, err := reloaded.GetBody()
	require.NoError(t, err)
	assert.Contains(t, body, "New body")
}

func TestUpdateNotePartial(t *testing.T) {
	tmp := withVault(t)

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\ntype: fleeting\n---\nOld body"), 0o644)
	require.NoError(t, err)

	input := UpdateNoteInput{
		Path:  path,
		Title: ptr.To("New Title"),
	}

	_, output, err := UpdateNote(t.Context(), nil, &input)
	require.NoError(t, err)
	assert.Equal(t, note.Type("fleeting"), output.Note.Frontmatter.Type)
	assert.Equal(t, "New Title", output.Note.Frontmatter.Title)
}

func TestUpdateNoteInvalidType(t *testing.T) {
	tmp := withVault(t)

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\ntype: fleeting\n---\nBody"), 0o644)
	require.NoError(t, err)

	input := UpdateNoteInput{
		Path: path,
		Type: ptr.To("not-a-type"),
	}

	_, _, err = UpdateNote(t.Context(), nil, &input)
	assert.Error(t, err)
}

func TestUpdateNoteNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := UpdateNote(t.Context(), nil, &UpdateNoteInput{Path: "note.md"})
	assert.Error(t, err)
}

func TestUpdateNoteNotFound(t *testing.T) {
	withVault(t)

	_, _, err := UpdateNote(t.Context(), nil, &UpdateNoteInput{Path: "does-not-exist.md"})
	assert.Error(t, err)
}
