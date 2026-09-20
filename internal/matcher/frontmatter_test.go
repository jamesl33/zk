package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrontmatterMatch(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			Title: "Test",
			Tags:  []string{"tag1"},
		},
	}

	m, err := Frontmatter("tag1", "", "", false)
	require.NoError(t, err)

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestFrontmatterMatchIgnoreCase(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			Title: "Test",
			Tags:  []string{"tag1"},
		},
	}

	m, err := Frontmatter("TAG1", "", "", true)
	require.NoError(t, err)

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestFrontmatterNoMatch(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			Title: "Test",
			Tags:  []string{"tag1"},
		},
	}

	m, err := Frontmatter("nonexistent", "", "", false)
	require.NoError(t, err)

	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}
