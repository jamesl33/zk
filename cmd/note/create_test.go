package note

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/jamesl33/zk/internal/note"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	tmp := t.TempDir()

	err := create("permanent", "Test Note", tmp)
	require.NoError(t, err)

	matches, err := filepath.Glob(filepath.Join(tmp, "*.md"))
	require.NoError(t, err)
	require.Len(t, matches, 1)

	assert.Regexp(t, regexp.MustCompile(`^\d{14}\.md$`), filepath.Base(matches[0]))

	n, err := note.New(matches[0])
	require.NoError(t, err)
	assert.Equal(t, note.Type("permanent"), n.Frontmatter.Type)
	assert.Equal(t, "Test Note", n.Frontmatter.Title)
}
