package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNameMatch(t *testing.T) {
	type test struct {
		name     string
		target   string
		notePath string
	}

	tests := []test{
		{
			name:     "match name",
			target:   "note",
			notePath: "note.md",
		},
		{
			name:     "match name with path",
			target:   "note",
			notePath: "subdir/note.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				m = Name(tt.target)
				n = &note.Note{Path: tt.notePath}
			)

			actual, err := m(n)
			require.NoError(t, err)
			assert.True(t, actual)
		})
	}
}

func TestNameNoMatch(t *testing.T) {
	var (
		m = Name("other")
		n = &note.Note{Path: "note.md"}
	)

	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}
