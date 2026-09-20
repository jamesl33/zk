package claude

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withTempDir(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	return tmp
}

func TestClaudeRun(t *testing.T) {
	tmp := withTempDir(t)

	var c Claude

	err := c.Run(t.Context())
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmp, "CLAUDE.md"))
	assert.FileExists(t, filepath.Join(tmp, ".mcp.json"))
	assert.FileExists(t, filepath.Join(tmp, ".claude", "settings.json"))

	entries, err := os.ReadDir(filepath.Join(tmp, ".claude", "skills"))
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestClaudeRunRemovesExisting(t *testing.T) {
	tmp := withTempDir(t)

	require.NoError(t, os.MkdirAll(filepath.Join(tmp, ".claude", "stale"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tmp, ".claude", "stale", "file.txt"), []byte("x"), 0o644))

	var c Claude

	err := c.Run(t.Context())
	require.NoError(t, err)

	assert.NoFileExists(t, filepath.Join(tmp, ".claude", "stale", "file.txt"))
}
