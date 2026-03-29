package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaggedMatch(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			Tags: []string{"tag1", "tag2"},
		},
	}

	m := Tagged("tag1")

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestTaggedNoMatch(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			Tags: []string{"tag1", "tag2"},
		},
	}

	m := Tagged("tag3")

	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}
