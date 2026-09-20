package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTitleMatch(t *testing.T) {
	m, err := Title("My Title", "", "", false)
	require.NoError(t, err)

	n := &note.Note{
		Frontmatter: note.Frontmatter{Title: "My Title"},
	}

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestTitleMatchIgnoreCase(t *testing.T) {
	m, err := Title("MY TITLE", "", "", true)
	require.NoError(t, err)

	n := &note.Note{
		Frontmatter: note.Frontmatter{Title: "My Title"},
	}

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestTitleNoMatch(t *testing.T) {
	m, err := Title("My Title", "", "", false)
	require.NoError(t, err)

	n := &note.Note{
		Frontmatter: note.Frontmatter{Title: "Other"},
	}

	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}
