package note

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
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

func TestCreateBibliographicRun(t *testing.T) {
	tmp := withVault(t)

	c := CreateBibliographic{CreateBibliographicOptions{Title: "My Book"}}

	err := c.Run(t.Context())
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(tmp, "5 Bibliography", "*.md"))
	require.NoError(t, err)
	require.Len(t, matches, 1)

	n, err := note.New(matches[0])
	require.NoError(t, err)
	assert.Equal(t, note.Type("bibliographic"), n.Frontmatter.Type)
	assert.Equal(t, "My Book", n.Frontmatter.Title)
}

func TestCreateFleetingRun(t *testing.T) {
	tmp := withVault(t)

	c := CreateFleeting{CreateFleetingOptions{Title: "Quick Thought"}}

	err := c.Run(t.Context())
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(tmp, "0 Inbox", "*.md"))
	require.NoError(t, err)
	require.Len(t, matches, 1)

	n, err := note.New(matches[0])
	require.NoError(t, err)
	assert.Equal(t, note.Type("fleeting"), n.Frontmatter.Type)
	assert.Equal(t, "Quick Thought", n.Frontmatter.Title)
}

func TestCreateIndexRun(t *testing.T) {
	tmp := withVault(t)

	c := CreateIndex{CreateIndexOptions{Title: "Area Index"}}

	err := c.Run(t.Context(), "1 Areas")
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(tmp, "1 Areas", "*.md"))
	require.NoError(t, err)
	require.Len(t, matches, 1)

	n, err := note.New(matches[0])
	require.NoError(t, err)
	assert.Equal(t, note.Type("index"), n.Frontmatter.Type)
	assert.Equal(t, "Area Index", n.Frontmatter.Title)
}

func TestCreateLiteratureRun(t *testing.T) {
	tmp := withVault(t)

	c := CreateLiterature{CreateLiteratureOptions{Title: "Literature Note"}}

	err := c.Run(t.Context(), "2 Literature")
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(tmp, "2 Literature", "*.md"))
	require.NoError(t, err)
	require.Len(t, matches, 1)

	n, err := note.New(matches[0])
	require.NoError(t, err)
	assert.Equal(t, note.Type("literature"), n.Frontmatter.Type)
	assert.Equal(t, "Literature Note", n.Frontmatter.Title)
}

func TestCreatePermanentRun(t *testing.T) {
	tmp := withVault(t)

	c := CreatePermanent{CreatePermanentOptions{Title: "Permanent Note"}}

	err := c.Run(t.Context(), "3 Permanent")
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(tmp, "3 Permanent", "*.md"))
	require.NoError(t, err)
	require.Len(t, matches, 1)

	n, err := note.New(matches[0])
	require.NoError(t, err)
	assert.Equal(t, note.Type("permanent"), n.Frontmatter.Type)
	assert.Equal(t, "Permanent Note", n.Frontmatter.Title)
}
