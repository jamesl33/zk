package matcher

import (
	"errors"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGtorMatch(t *testing.T) {
	type test struct {
		glob  string
		input string
	}

	tests := []test{
		{glob: "*.md", input: "note.md"},
		{glob: "note?.md", input: "note1.md"},
		{glob: "note[12].md", input: "note1.md"},
		{glob: "note[12].md", input: "note2.md"},
		{glob: "note[!1].md", input: "note2.md"},
	}

	for _, tt := range tests {
		t.Run(tt.glob+"_"+tt.input, func(t *testing.T) {
			pattern := gtor(tt.glob)
			assert.Regexp(t, "^"+pattern+"$", tt.input)
		})
	}
}

func TestGtorNoMatch(t *testing.T) {
	type test struct {
		glob  string
		input string
	}

	tests := []test{
		{glob: "*.md", input: "note.txt"},
		{glob: "note?.md", input: "note.md"},
		{glob: "note[12].md", input: "note3.md"},
		{glob: "note[!1].md", input: "note1.md"},
	}

	for _, tt := range tests {
		t.Run(tt.glob+"_"+tt.input, func(t *testing.T) {
			pattern := gtor(tt.glob)
			assert.NotRegexp(t, "^"+pattern+"$", tt.input)
		})
	}
}

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

func TestTextError(t *testing.T) {
	extract := func(n *note.Note) (string, error) { return "", errors.New("error") }

	m, err := text("pattern", "", "", extract)
	require.NoError(t, err)

	_, err = m(&note.Note{})
	assert.Error(t, err)
}
