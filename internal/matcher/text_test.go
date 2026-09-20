package matcher

import (
	"errors"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTextMatch(t *testing.T) {
	extract := func(n *note.Note) (string, error) { return n.Frontmatter.Title, nil }

	type test struct {
		name  string
		fixed string
		glob  string
		regex string
		title string
	}

	tests := []test{
		{
			name:  "fixed match",
			fixed: "Note",
			title: "My Note",
		},
		{
			name:  "glob match",
			glob:  "*.md",
			title: "note.md",
		},
		{
			name:  "regex match",
			regex: "Note [0-9]",
			title: "Note 1 ",
		},
		{
			name:  "multiple matches (OR)",
			fixed: "Note",
			glob:  "Other*",
			title: "Other Note",
		},
		{
			name:  "empty patterns match anything",
			title: "Anything",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := text(tt.fixed, tt.glob, tt.regex, extract)
			require.NoError(t, err)

			n := &note.Note{Frontmatter: note.Frontmatter{Title: tt.title}}
			actual, err := m(n)
			require.NoError(t, err)
			assert.True(t, actual)
		})
	}
}

func TestTextNoMatch(t *testing.T) {
	extract := func(n *note.Note) (string, error) { return n.Frontmatter.Title, nil }

	m, err := text("Other", "", "", extract)
	require.NoError(t, err)

	n := &note.Note{Frontmatter: note.Frontmatter{Title: "My Note"}}
	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}

func TestTextMatchSmartCaseInsensitive(t *testing.T) {
	extract := func(n *note.Note) (string, error) { return n.Frontmatter.Title, nil }

	m, err := text("my note", "", "", extract)
	require.NoError(t, err)

	n := &note.Note{Frontmatter: note.Frontmatter{Title: "MY NOTE"}}
	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestTextNoMatchSmartCaseSensitive(t *testing.T) {
	extract := func(n *note.Note) (string, error) { return n.Frontmatter.Title, nil }

	m, err := text("MY NOTE", "", "", extract)
	require.NoError(t, err)

	n := &note.Note{Frontmatter: note.Frontmatter{Title: "my note"}}
	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}

func TestTextError(t *testing.T) {
	extract := func(n *note.Note) (string, error) { return "", errors.New("error") }

	m, err := text("pattern", "", "", extract)
	require.NoError(t, err)

	_, err = m(&note.Note{})
	assert.Error(t, err)
}

func TestCaseInsensitive(t *testing.T) {
	assert.True(t, insensitive())
	assert.True(t, insensitive("lower", "", ""))
	assert.False(t, insensitive("Upper", "", ""))
	assert.False(t, insensitive("", "Upper", ""))
	assert.False(t, insensitive("", "", "Upper"))
}
