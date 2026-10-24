package memory

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAutoMemoryEnabled(t *testing.T) {
	home, cwd := t.TempDir(), t.TempDir()

	t.Setenv("HOME", home)
	t.Setenv("CLAUDE_CODE_DISABLE_AUTO_MEMORY", "")

	require.True(t, AutoMemoryEnabled(cwd))

	require.NoError(t, os.MkdirAll(filepath.Join(home, ".claude"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(home, ".claude", "settings.json"), []byte(`{"autoMemoryEnabled": false}`), 0o644))
	require.False(t, AutoMemoryEnabled(cwd))
}

func TestAutoMemoryEnabledEnv(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("CLAUDE_CODE_DISABLE_AUTO_MEMORY", "1")

	require.False(t, AutoMemoryEnabled(t.TempDir()))
}
