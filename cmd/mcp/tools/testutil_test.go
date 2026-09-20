package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// withVault creates a temporary vault directory, chdirs into it for the
// duration of the test, and returns its path.
func withVault(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	return tmp
}

// withEmptyDir chdirs into a temporary directory (without a '.zk' vault
// marker) for the duration of the test.
func withEmptyDir(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	return tmp
}
