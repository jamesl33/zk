package gemini

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

func TestGeminiRun(t *testing.T) {
	tmp := withTempDir(t)

	var g Gemini

	err := g.Run(t.Context())
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(tmp, "GEMINI.md"))
	assert.FileExists(t, filepath.Join(tmp, ".gemini", "settings.json"))
	assert.FileExists(t, filepath.Join(tmp, ".gemini", "policies", "zk.toml"))

	entries, err := os.ReadDir(filepath.Join(tmp, ".gemini", "skills"))
	require.NoError(t, err)
	assert.NotEmpty(t, entries)
}

func TestGeminiRunRemovesExisting(t *testing.T) {
	tmp := withTempDir(t)

	require.NoError(t, os.MkdirAll(filepath.Join(tmp, ".gemini", "stale"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(tmp, ".gemini", "stale", "file.txt"), []byte("x"), 0o644))

	var g Gemini

	err := g.Run(t.Context())
	require.NoError(t, err)

	assert.NoFileExists(t, filepath.Join(tmp, ".gemini", "stale", "file.txt"))
}
