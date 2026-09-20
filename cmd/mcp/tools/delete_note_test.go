package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteNote(t *testing.T) {
	tmp := withVault(t)

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\ntype: fleeting\n---\nBody"), 0o644)
	require.NoError(t, err)

	_, output, err := DeleteNote(t.Context(), nil, &DeleteNoteInput{Path: path})
	require.NoError(t, err)
	require.NotNil(t, output)

	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))
}

func TestDeleteNoteNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := DeleteNote(t.Context(), nil, &DeleteNoteInput{Path: "note.md"})
	assert.Error(t, err)
}

func TestDeleteNoteNotFound(t *testing.T) {
	withVault(t)

	_, _, err := DeleteNote(t.Context(), nil, &DeleteNoteInput{Path: "does-not-exist.md"})
	assert.Error(t, err)
}
