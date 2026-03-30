package notes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	m := func(n *note.Note) (bool, error) {
		return n.Frontmatter.Title == "Note 2", nil
	}

	var actual []*note.Note

	err = Search(t.Context(), tmp, m, func(n *note.Note) {
		actual = append(actual, n)
	})
	require.NoError(t, err)
	assert.Len(t, actual, 1)
	assert.Equal(t, "Note 2", actual[0].Frontmatter.Title)
}

func TestSearchNoMatches(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	m := func(n *note.Note) (bool, error) {
		return false, nil
	}

	var actual []*note.Note

	err = Search(t.Context(), tmp, m, func(n *note.Note) {
		actual = append(actual, n)
	})
	require.NoError(t, err)
	assert.Empty(t, actual)
}

func TestSearchError(t *testing.T) {
	tmp := t.TempDir()

	// Path does not exist
	err := Search(t.Context(), filepath.Join(tmp, "non-existent"), matcher.Any(), func(n *note.Note) {})
	assert.Error(t, err)
}

func TestSearchMatcherError(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	m := func(n *note.Note) (bool, error) {
		return false, assert.AnError
	}

	err = Search(t.Context(), tmp, m, func(n *note.Note) {})
	assert.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
