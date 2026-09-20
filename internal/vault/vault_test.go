package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootFindsZkMarker(t *testing.T) {
	root := t.TempDir()

	require.NoError(t, os.Mkdir(filepath.Join(root, ".zk"), 0o755))

	sub := filepath.Join(root, "a", "b")
	require.NoError(t, os.MkdirAll(sub, 0o755))

	found, err := Root(sub)
	require.NoError(t, err)
	assert.Equal(t, root, found)
}

func TestRootNotFound(t *testing.T) {
	_, err := Root(t.TempDir())
	assert.ErrorIs(t, err, ErrNotFound)
}
