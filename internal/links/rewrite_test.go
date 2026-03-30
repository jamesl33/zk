package links

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRewrite(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmp)
	require.NoError(t, err)
	defer os.Chdir(cwd)

	err = os.WriteFile(
		filepath.Join(tmp, "20251129123456.md"),
		[]byte("---\ntype: permanent\ntitle: Linked Note Title\n---\nBody"),
		0o644,
	)
	require.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(tmp, "main.md"),
		[]byte("---\ntype: permanent\ntitle: Main Note\n---\nLink to [[20251129123456]] and [[20251129123457]]"),
		0o644,
	)
	require.NoError(t, err)

	n, err := note.New(filepath.Join(tmp, "main.md"))
	require.NoError(t, err)

	err = Rewrite(context.Background(), n)
	require.NoError(t, err)

	body, err := n.GetBody()
	require.NoError(t, err)
	assert.Equal(t, "Link to [[20251129123456|Linked Note Title]] and [[20251129123457]]", body)
}

func TestRewriteUpdate(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmp)
	require.NoError(t, err)
	defer os.Chdir(cwd)

	err = os.WriteFile(
		filepath.Join(tmp, "20251129123456.md"),
		[]byte("---\ntype: permanent\ntitle: New Title\n---\nBody"),
		0o644,
	)
	require.NoError(t, err)

	err = os.WriteFile(
		filepath.Join(tmp, "main.md"),
		[]byte("---\ntype: permanent\ntitle: Main Note\n---\nLink to [[20251129123456|Old Title]]"),
		0o644,
	)
	require.NoError(t, err)

	n, err := note.New(filepath.Join(tmp, "main.md"))
	require.NoError(t, err)

	err = Rewrite(context.Background(), n)
	require.NoError(t, err)

	body, err := n.GetBody()
	require.NoError(t, err)
	assert.Equal(t, "Link to [[20251129123456|New Title]]", body)
}
