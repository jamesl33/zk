package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyMatch(t *testing.T) {
	n := &note.Note{}
	n.SetBody("This is the body")

	m, err := Body("the body", "", "")
	require.NoError(t, err)

	actual, err := m(n)
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestBodyNoMatch(t *testing.T) {
	n := &note.Note{}
	n.SetBody("Other content")

	m, err := Body("the body", "", "")
	require.NoError(t, err)

	actual, err := m(n)
	require.NoError(t, err)
	assert.False(t, actual)
}
