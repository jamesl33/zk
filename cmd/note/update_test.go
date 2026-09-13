package note

import (
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// withStdin replaces os.Stdin with the given content for the duration of the
// test.
func withStdin(t *testing.T, content string) {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	_, err = w.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	original := os.Stdin
	os.Stdin = r

	t.Cleanup(func() { os.Stdin = original })
}

func TestUpdatePathFromArgs(t *testing.T) {
	var u Update

	path, err := u.path([]string{"note.md"})
	require.NoError(t, err)
	assert.Equal(t, "note.md", path)
}

func TestUpdatePathFromStdin(t *testing.T) {
	withStdin(t, "note.md\n")

	var u Update

	path, err := u.path(nil)
	require.NoError(t, err)
	assert.Equal(t, "note.md", path)
}

func TestUpdatePathFromStdinNoTrailingNewline(t *testing.T) {
	withStdin(t, "note.md")

	var u Update

	path, err := u.path([]string{"-"})
	require.NoError(t, err)
	assert.Equal(t, "note.md", path)
}

func TestUpdatePathFromStdinEmpty(t *testing.T) {
	withStdin(t, "")

	var u Update

	path, err := u.path(nil)
	require.True(t, errors.Is(err, io.EOF))
	assert.Empty(t, path)
}
