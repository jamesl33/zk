package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnterVault(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	t.Cleanup(func() { require.NoError(t, os.Chdir(wd)) })

	require.NoError(t, enterVault(""))

	got, err := os.Getwd()
	require.NoError(t, err)
	require.Equal(t, wd, got)

	tmp, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)

	require.NoError(t, enterVault(tmp))

	got, err = os.Getwd()
	require.NoError(t, err)
	require.Equal(t, tmp, got)

	require.Error(t, enterVault(filepath.Join(tmp, "missing")))
}
