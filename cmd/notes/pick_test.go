package notes

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPickItem(t *testing.T) {
	var p Pick

	out := captureStdout(t, func() {
		buf := bytes.NewBufferString("col1\x01col2\x01col3\x01/path/to/note.md")
		err := p.item(*buf)
		require.NoError(t, err)
	})

	assert.Equal(t, "/path/to/note.md", out)
}

func TestPickItemEmpty(t *testing.T) {
	var p Pick

	out := captureStdout(t, func() {
		buf := bytes.NewBufferString("")
		err := p.item(*buf)
		require.NoError(t, err)
	})

	assert.Empty(t, out)
}

func TestPickItemNewlineOnly(t *testing.T) {
	var p Pick

	out := captureStdout(t, func() {
		buf := bytes.NewBufferString("\n")
		err := p.item(*buf)
		require.NoError(t, err)
	})

	assert.Empty(t, out)
}
