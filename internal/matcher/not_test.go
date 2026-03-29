package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotMatch(t *testing.T) {
	m := Not(nmatch)

	actual, err := m(&note.Note{})
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestNotNoMatch(t *testing.T) {
	m := Not(match)

	actual, err := m(&note.Note{})
	require.NoError(t, err)
	assert.False(t, actual)
}

func TestNotError(t *testing.T) {
	m := Not(ematch)

	_, err := m(&note.Note{})
	assert.Error(t, err)
}
