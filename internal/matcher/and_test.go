package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	match  = func(n *note.Note) (bool, error) { return true, nil }
	nmatch = func(n *note.Note) (bool, error) { return false, nil }
	ematch = func(n *note.Note) (bool, error) { return false, assert.AnError }
)

func TestAndMatch(t *testing.T) {
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
			name:     "all true matchers",
			matchers: []Matcher{match, match},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := And(tt.matchers...)
			actual, err := m(&note.Note{})
			require.NoError(t, err)
			assert.True(t, actual)
		})
	}
}

func TestAndNoMatch(t *testing.T) {
	type test struct {
		name     string
		matchers []Matcher
	}

	tests := []test{
		{
			name:     "one false matcher",
			matchers: []Matcher{match, nmatch},
		},
		{
			name:     "all false matchers",
			matchers: []Matcher{nmatch, nmatch},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := And(tt.matchers...)

			actual, err := m(&note.Note{})
			require.NoError(t, err)
			assert.False(t, actual)
		})
	}
}

func TestAndError(t *testing.T) {
	m := And(match, ematch)

	_, err := m(&note.Note{})
	assert.ErrorIs(t, err, assert.AnError)
}
