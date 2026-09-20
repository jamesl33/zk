package links

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestRewriteRewrite(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\n[[20060102150405]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	n, err := note.New(filepath.Join(tmp, "20060102150404.md"))
	require.NoError(t, err)

	var r Rewrite

	err = r.rewrite(t.Context(), n)
	require.NoError(t, err)

	reloaded, err := note.New(filepath.Join(tmp, "20060102150404.md"))
	require.NoError(t, err)

	body, err := reloaded.GetBody()
	require.NoError(t, err)
	assert.Contains(t, body, "[[20060102150405|Note 2]]")
}

func TestRewriteRun(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\n[[20060102150405]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	var r Rewrite

	err = r.Run(t.Context(), nil)
	require.NoError(t, err)

	reloaded, err := note.New(filepath.Join(tmp, "20060102150404.md"))
	require.NoError(t, err)

	body, err := reloaded.GetBody()
	require.NoError(t, err)
	assert.Contains(t, body, "[[20060102150405|Note 2]]")
}

func TestRewriteRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var r Rewrite

	err = r.Run(t.Context(), nil)
	assert.Error(t, err)
}

func TestRewriteRunWithPathArg(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\n[[20060102150405]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	var r Rewrite

	err = r.Run(t.Context(), []string{"."})
	assert.NoError(t, err)
}
