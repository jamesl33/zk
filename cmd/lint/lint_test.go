package lint

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

func TestLintRun(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\ntype: permanent\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	var l Lint

	err = l.Run(t.Context(), nil)
	assert.NoError(t, err)
}

func TestLintRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var l Lint

	err = l.Run(t.Context(), nil)
	assert.Error(t, err)
}

func TestLintRunWithPathArg(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\ntype: permanent\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	var l Lint

	err = l.Run(t.Context(), []string{"."})
	assert.NoError(t, err)
}

func TestLintRunOrphanNote(t *testing.T) {
	tmp := withVault(t)

	// A permanent note with no links is flagged as an orphan (line/column both 0).
	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\ntype: permanent\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	var l Lint

	out := captureStdout(t, func() {
		err = l.Run(t.Context(), nil)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "orphan-note")
	assert.NotContains(t, out, ":0:0:")
}

func TestLintRunBrokenLink(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\n[[20060102150405]]"), 0o644)
	require.NoError(t, err)

	var l Lint

	out := captureStdout(t, func() {
		err = l.Run(t.Context(), nil)
		require.NoError(t, err)
	})

	assert.Contains(t, out, "linkcheck")
}

func TestLintRunArchivesFlag(t *testing.T) {
	tmp := withVault(t)

	require.NoError(t, os.Mkdir(filepath.Join(tmp, "4 Archives"), 0o755))

	err := os.WriteFile(
		filepath.Join(tmp, "4 Archives", "note1.md"),
		[]byte("---\ntitle: Note 1\ntype: permanent\n---\nBody 1"),
		0o644,
	)
	require.NoError(t, err)

	l := Lint{LintOptions{Archives: true}}

	err = l.Run(t.Context(), nil)
	assert.NoError(t, err)
}
