package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadNote(t *testing.T) {
	tmp := withVault(t)

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	_, output, err := ReadNote(t.Context(), nil, &ReadNoteInput{Path: path})
	require.NoError(t, err)
	assert.Equal(t, "Note 1", output.Note.Frontmatter.Title)
	assert.Equal(t, "Body 1", output.Body)
}

func TestReadNoteNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := ReadNote(t.Context(), nil, &ReadNoteInput{Path: "note.md"})
	assert.Error(t, err)
}

func TestReadNoteNotFound(t *testing.T) {
	withVault(t)

	_, _, err := ReadNote(t.Context(), nil, &ReadNoteInput{Path: "does-not-exist.md"})
	assert.Error(t, err)
}
