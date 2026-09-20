package note

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var u Update

	err = u.Run(t.Context(), []string{"note.md"})
	assert.Error(t, err)
}

func TestUpdateRunEOFExitsCleanly(t *testing.T) {
	tmp := withVault(t)
	_ = tmp

	withStdin(t, "")

	var u Update

	err := u.Run(t.Context(), nil)
	assert.NoError(t, err)
}

func TestUpdateRunNoteNotFound(t *testing.T) {
	_ = withVault(t)

	var u Update

	err := u.Run(t.Context(), []string{"does-not-exist.md"})
	assert.Error(t, err)
}

func TestUpdateRunEditError(t *testing.T) {
	tmp := withVault(t)

	require.NoError(t, os.Unsetenv("EDITOR"))

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	var u Update

	err = u.Run(t.Context(), []string{path})
	assert.Error(t, err)
}
