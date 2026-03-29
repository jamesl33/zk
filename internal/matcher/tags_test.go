package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagsMatch(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			Tags: []string{"tag1", "tag2"},
		},
	}

	type test struct {
		name    string
		include []string
		exclude []string
	}

	tests := []test{
		{
			name:    "include match",
			include: []string{"tag1"},
		},
		{
			name:    "include match AND exclude no match",
			include: []string{"tag1"},
			exclude: []string{"tag3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Tags(tt.include, tt.exclude)
			require.NoError(t, err)

			actual, err := m(n)
			require.NoError(t, err)
			assert.True(t, actual)
		})
	}
}

func TestTagsNoMatch(t *testing.T) {
	n := &note.Note{
		Frontmatter: note.Frontmatter{
			Tags: []string{"tag1", "tag2"},
		},
	}

	type test struct {
		name    string
		include []string
		exclude []string
	}

	tests := []test{
		{
			name:    "include no match",
			include: []string{"tag3"},
		},
		{
			name:    "exclude match",
			exclude: []string{"tag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Tags(tt.include, tt.exclude)
			require.NoError(t, err)

			actual, err := m(n)
			require.NoError(t, err)
			assert.False(t, actual)
		})
	}
}
