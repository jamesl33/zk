package note

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummarizeRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var s Summarize

	err = s.Run(t.Context(), "note.md")
	assert.Error(t, err)
}

func TestSummarizeRunNoteNotFound(t *testing.T) {
	_ = withVault(t)

	var s Summarize

	err := s.Run(t.Context(), "does-not-exist.md")
	assert.Error(t, err)
}

func TestSummarizeRunEmptyBody(t *testing.T) {
	tmp := withVault(t)

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\n---\n"), 0o644)
	require.NoError(t, err)

	var s Summarize

	// No AI client is reachable in this test, so a non-empty body would fail;
	// an empty body must return early without calling out.
	err = s.Run(t.Context(), path)
	assert.NoError(t, err)
}
