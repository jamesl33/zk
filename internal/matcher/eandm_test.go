package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEandmMatch(t *testing.T) {
	m := eandm(
		func(n *note.Note) (string, error) { return "hello", nil },
		func(text string) bool { return text == "hello" },
	)
	actual, err := m(&note.Note{})
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestEandmNoMatch(t *testing.T) {
	m := eandm(
		func(n *note.Note) (string, error) { return "hello", nil },
		func(text string) bool { return text == "world" },
	)
	actual, err := m(&note.Note{})
	require.NoError(t, err)
	assert.False(t, actual)
}

func TestEandmError(t *testing.T) {
	m := eandm(
		func(n *note.Note) (string, error) { return "", assert.AnError },
		func(text string) bool { return true },
	)
	_, err := m(&note.Note{})
	assert.ErrorIs(t, err, assert.AnError)
}
