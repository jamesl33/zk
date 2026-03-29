package matcher

import (
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnyMatch(t *testing.T) {
	m := Any()

	actual, err := m(&note.Note{})
	require.NoError(t, err)
	assert.True(t, actual)
}
