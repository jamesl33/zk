package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntireMatch(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{Title: "Title"},
	}
	n.SetBody("Body content")

	// Match title
	m, err := Entire("Title", "", "")
	require.NoError(t, err)

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)

	// Match body
	m, err = Entire("content", "", "")
	require.NoError(t, err)

	actual, err = m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestEntireNoMatch(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{Title: "Title"},
	}

	n.SetBody("Body content")

	m, err := Entire("nothing", "", "")
	require.NoError(t, err)

	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}
