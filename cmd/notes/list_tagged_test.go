package notes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListTaggedRun(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\ntags:\n  - hiking\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Note 2\ntags:\n  - cooking\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	l := ListTagged{ListTaggedOptions{With: []string{"hiking"}}}

	out := captureStdout(t, func() {
		err = l.Run(t.Context(), nil)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "Note 1")
	assert.NotContains(t, out, "Note 2")
}

func TestListTaggedRunWithout(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\ntags:\n  - hiking\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Note 2\ntags:\n  - cooking\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	l := ListTagged{ListTaggedOptions{Without: []string{"hiking"}}}

	out := captureStdout(t, func() {
		err = l.Run(t.Context(), []string{"."})
		require.NoError(t, err)
	})

	assert.Contains(t, out, "Note 2")
	assert.NotContains(t, out, "Note 1")
}

func TestListTaggedRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var l ListTagged

	err = l.Run(t.Context(), nil)
	assert.Error(t, err)
}
