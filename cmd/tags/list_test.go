package tags

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureStdout redirects os.Stdout for the duration of fn, returning what was written.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	original := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = original })

	fn()

	require.NoError(t, w.Close())
	os.Stdout = original

	out, err := io.ReadAll(r)
	require.NoError(t, err)

	return string(out)
}

func TestListRun(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	err = os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\ntags:\n  - b\n  - a\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Note 2\ntags:\n  - a\n  - c\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	var l List

	out := captureStdout(t, func() {
		err = l.Run(t.Context(), nil)
		require.NoError(t, err)
	})

	assert.Equal(t, "a\nb\nc\n", out)
}

func TestListRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var l List

	err = l.Run(t.Context(), nil)
	assert.Error(t, err)
}

func TestListRunWithPathArg(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	err = os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\ntags:\n  - a\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	var l List

	out := captureStdout(t, func() {
		err = l.Run(t.Context(), []string{filepath.Join(tmp, "note1.md")})
		require.NoError(t, err)
	})

	assert.Equal(t, "a\n", out)
}
