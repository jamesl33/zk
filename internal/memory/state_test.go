package memory

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRecalledLifecycle(t *testing.T) {
	vault := withVault(t)

	recalled, err := Recalled(vault, "s")
	require.NoError(t, err)
	require.Empty(t, recalled)

	require.NoError(t, MarkRecalled(vault, "s", "a.md", "b.md"))

	recalled, err = Recalled(vault, "s")
	require.NoError(t, err)
	require.Equal(t, map[string]bool{"a.md": true, "b.md": true}, recalled)

	require.NoError(t, ResetRecalled(vault, "s"))
	require.NoError(t, ResetRecalled(vault, "s")) // Nothing to reset is fine

	recalled, err = Recalled(vault, "s")
	require.NoError(t, err)
	require.Empty(t, recalled)
}

func TestSessionIDCannotEscape(t *testing.T) {
	vault := withVault(t)

	require.NoError(t, SetOffset(vault, "../../escape", 1))
	require.FileExists(t, filepath.Join(stateDir(vault), "escape.offset"))
}

func TestOffset(t *testing.T) {
	vault := withVault(t)

	require.Zero(t, Offset(vault, "s"))
	require.NoError(t, SetOffset(vault, "s", 42))
	require.Equal(t, 42, Offset(vault, "s"))
}

func TestPruneKeepsRecentAndLog(t *testing.T) {
	vault := withVault(t)

	require.NoError(t, SetOffset(vault, "old", 1))
	require.NoError(t, SetOffset(vault, "new", 1))
	logPath := filepath.Join(stateDir(vault), "hooks.log")
	require.NoError(t, os.WriteFile(logPath, []byte("x"), 0o644))

	old := time.Now().Add(-2 * stateMaxAge)
	require.NoError(t, os.Chtimes(sessionFile(vault, "old", "offset"), old, old))
	require.NoError(t, os.Chtimes(logPath, old, old))

	require.NoError(t, Prune(vault))

	require.NoFileExists(t, sessionFile(vault, "old", "offset"))
	require.FileExists(t, sessionFile(vault, "new", "offset"))
	require.FileExists(t, logPath)
}

func TestLogfStartsAfreshWhenLarge(t *testing.T) {
	vault := withVault(t)

	path := filepath.Join(stateDir(vault), "hooks.log")

	require.NoError(t, os.WriteFile(path, make([]byte, logMaxSize+1), 0o644))

	Logf(vault, "hello %d", 1)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Contains(t, string(data), "hello 1")
	require.Less(t, len(data), 100)
}
