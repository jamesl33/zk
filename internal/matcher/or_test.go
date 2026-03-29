package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrMatch(t *testing.T) {
	type test struct {
		name     string
		matchers []Matcher
	}

	tests := []test{
		{
			name:     "empty matchers",
			matchers: []Matcher{},
		},
		{
			name:     "one true matcher",
			matchers: []Matcher{nmatch, match},
		},
		{
			name:     "all true matchers",
			matchers: []Matcher{match, match},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Or(tt.matchers...)

			actual, err := m(&note.Note{})
			require.NoError(t, err)
			assert.True(t, actual)
		})
	}
}

func TestOrNoMatch(t *testing.T) {
	m := Or(nmatch, nmatch)

	actual, err := m(&note.Note{})
	require.NoError(t, err)
	assert.False(t, actual)
}

func TestOrError(t *testing.T) {
	m := Or(nmatch, ematch)

	_, err := m(&note.Note{})
	assert.Error(t, err)
}
