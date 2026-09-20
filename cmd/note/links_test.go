package note

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/note"
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

func TestLinksToDisabled(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	n, err := note.New(filepath.Join(tmp, "note1.md"))
	require.NoError(t, err)

	l := Links{}

	out := captureStdout(t, func() {
		err = l.to(t.Context(), n)
		require.NoError(t, err)
	})

	assert.Empty(t, out)
}

func TestLinksToEnabled(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	err = os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\n[[20060102150404]]"), 0o644)
	require.NoError(t, err)

	n, err := note.New(filepath.Join(tmp, "20060102150404.md"))
	require.NoError(t, err)

	l := Links{LinksOptions{To: true}}

	out := captureStdout(t, func() {
		err = l.to(t.Context(), n)
		require.NoError(t, err)
	})

	assert.NotEmpty(t, out)
}

func TestLinksFromDisabled(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	n, err := note.New(filepath.Join(tmp, "note1.md"))
	require.NoError(t, err)

	l := Links{}

	out := captureStdout(t, func() {
		err = l.from(t.Context(), n)
		require.NoError(t, err)
	})

	assert.Empty(t, out)
}

func TestLinksFromEnabled(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	err = os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\n[[20060102150405]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	n, err := note.New(filepath.Join(tmp, "20060102150404.md"))
	require.NoError(t, err)

	l := Links{LinksOptions{From: true}}

	out := captureStdout(t, func() {
		err = l.from(t.Context(), n)
		require.NoError(t, err)
	})

	assert.NotEmpty(t, out)
}

func TestLinksRunDefaultsBoth(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	err = os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	var l Links

	out := captureStdout(t, func() {
		err = l.Run(t.Context(), "20060102150404.md")
		require.NoError(t, err)
	})

	assert.Empty(t, out)
	assert.True(t, l.To)
	assert.True(t, l.From)
}

func TestLinksRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var l Links

	err = l.Run(t.Context(), "note.md")
	assert.Error(t, err)
}

func TestLinksRunNoteNotFound(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	var l Links

	err = l.Run(t.Context(), "does-not-exist.md")
	assert.Error(t, err)
}
