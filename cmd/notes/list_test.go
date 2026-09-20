package notes

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// withVault creates a temporary vault directory, chdirs into it for the
// duration of the test, and returns its path.
func withVault(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	return tmp
}

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
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Hiking Trip\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Cooking Recipe\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	l := List{ListOptions{Fixed: "Hiking"}}

	out := captureStdout(t, func() {
		err = l.Run(t.Context(), nil)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "Hiking Trip")
	assert.NotContains(t, out, "Cooking Recipe")
}

func TestListRunNoFilter(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	var l List

	out := captureStdout(t, func() {
		err = l.Run(t.Context(), []string{"."})
		require.NoError(t, err)
	})

	assert.Contains(t, out, "Note 1")
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

func TestListRunInvalidRegex(t *testing.T) {
	withVault(t)

	l := List{ListOptions{Regex: "("}}

	err := l.Run(t.Context(), nil)
	assert.Error(t, err)
}
