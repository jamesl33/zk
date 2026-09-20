package tags

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteUpdate(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\ntags:\n  - keep\n  - remove\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	n, err := note.New(path)
	require.NoError(t, err)

	var d Delete

	err = d.update(n, "remove")
	require.NoError(t, err)
	assert.Equal(t, []string{"keep"}, n.Frontmatter.Tags)

	reloaded, err := note.New(path)
	require.NoError(t, err)
	assert.Equal(t, []string{"keep"}, reloaded.Frontmatter.Tags)
}

func TestDeleteUpdateNotPresent(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\ntags:\n  - keep\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	n, err := note.New(path)
	require.NoError(t, err)

	var d Delete

	err = d.update(n, "remove")
	require.NoError(t, err)
	assert.Equal(t, []string{"keep"}, n.Frontmatter.Tags)
}

func TestDeleteRun(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	err = os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\ntags:\n  - keep\n  - remove\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Note 2\ntags:\n  - other\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	var d Delete

	err = d.Run(t.Context(), "remove")
	require.NoError(t, err)

	n1, err := note.New(filepath.Join(tmp, "note1.md"))
	require.NoError(t, err)
	assert.Equal(t, []string{"keep"}, n1.Frontmatter.Tags)

	n2, err := note.New(filepath.Join(tmp, "note2.md"))
	require.NoError(t, err)
	assert.Equal(t, []string{"other"}, n2.Frontmatter.Tags)
}

func TestDeleteRunNoVault(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	var d Delete

	err = d.Run(t.Context(), "remove")
	assert.Error(t, err)
}
