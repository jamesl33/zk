package vault

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRun(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "vault")

	var i Init

	require.NoError(t, i.Run(tmp))

	for _, dir := range append([]string{".zk"}, directories...) {
		assert.DirExists(t, filepath.Join(tmp, dir))
	}

	// Running again must leave existing content untouched
	note := filepath.Join(tmp, "0 Inbox", "note.md")
	require.NoError(t, os.WriteFile(note, []byte("keep"), 0o644))

	require.NoError(t, i.Run(tmp))

	data, err := os.ReadFile(note)
	require.NoError(t, err)
	assert.Equal(t, "keep", string(data))
}
