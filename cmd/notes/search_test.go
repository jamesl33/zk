package notes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchRun(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody about hiking"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Note 2\n---\nBody about cooking"), 0o644)
	require.NoError(t, err)

	s := Search{SearchOptions{Fixed: "hiking"}}

	out := captureStdout(t, func() {
		err = s.Run(t.Context(), nil)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "Note 1")
	assert.NotContains(t, out, "Note 2")
}

func TestSearchRunWithPathArg(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	var s Search

	out := captureStdout(t, func() {
		err = s.Run(t.Context(), []string{"."})
		require.NoError(t, err)
	})

	assert.Contains(t, out, "Note 1")
}

func TestSearchRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var s Search

	err = s.Run(t.Context(), nil)
	assert.Error(t, err)
}

func TestSearchRunInvalidRegex(t *testing.T) {
	withVault(t)

	s := Search{SearchOptions{Regex: "("}}

	err := s.Run(t.Context(), nil)
	assert.Error(t, err)
}
