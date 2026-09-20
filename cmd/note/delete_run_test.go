package note

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var d Delete

	err = d.Run(t.Context(), []string{"note.md"})
	assert.Error(t, err)
}

func TestDeleteRunEOFExitsCleanly(t *testing.T) {
	tmp := withVault(t)
	_ = tmp

	withStdin(t, "")

	var d Delete

	err := d.Run(t.Context(), nil)
	assert.NoError(t, err)
}

func TestDeleteRunNoteNotFound(t *testing.T) {
	_ = withVault(t)

	var d Delete

	err := d.Run(t.Context(), []string{"does-not-exist.md"})
	assert.Error(t, err)
}

func TestDeleteRunRemovesNote(t *testing.T) {
	tmp := withVault(t)

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	var d Delete

	err = d.Run(t.Context(), []string{path})
	require.NoError(t, err)

	_, err = os.Stat(path)
	assert.True(t, os.IsNotExist(err))
}
