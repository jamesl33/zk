package note

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeletePathFromArgs(t *testing.T) {
	var d Delete

	path, err := d.path([]string{"note.md"})
	require.NoError(t, err)
	assert.Equal(t, "note.md", path)
}

func TestDeletePathFromStdin(t *testing.T) {
	withStdin(t, "note.md\n")

	var d Delete

	path, err := d.path(nil)
	require.NoError(t, err)
	assert.Equal(t, "note.md", path)
}

func TestDeletePathFromStdinNoTrailingNewline(t *testing.T) {
	withStdin(t, "note.md")

	var d Delete

	path, err := d.path([]string{"-"})
	require.NoError(t, err)
	assert.Equal(t, "note.md", path)
}

func TestDeletePathFromStdinEmpty(t *testing.T) {
	withStdin(t, "")

	var d Delete

	path, err := d.path(nil)
	require.True(t, errors.Is(err, io.EOF))
	assert.Empty(t, path)
}
