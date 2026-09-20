package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathMatch(t *testing.T) {
	m, err := Path("note.md", "", "")
	require.NoError(t, err)

	n := &note.Note{
		Path: "note.md",
	}

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestPathMatchSmartCaseInsensitive(t *testing.T) {
	m, err := Path("note.md", "", "")
	require.NoError(t, err)

	n := &note.Note{
		Path: "NOTE.MD",
	}

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestPathMatchSmartCaseSensitive(t *testing.T) {
	m, err := Path("NOTE.MD", "", "")
	require.NoError(t, err)

	n := &note.Note{
		Path: "NOTE.MD",
	}

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestPathNoMatch(t *testing.T) {
	m, err := Path("note.md", "", "")
	require.NoError(t, err)

	n := &note.Note{
		Path: "other.md",
	}

	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}

func TestPathNoMatchSmartCaseSensitive(t *testing.T) {
	m, err := Path("NOTE.MD", "", "")
	require.NoError(t, err)

	n := &note.Note{
		Path: "note.md",
	}

	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}
