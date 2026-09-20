package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindNotesLinkedFrom(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\n[[20060102150405]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	path := filepath.Join(tmp, "20060102150404.md")

	_, output, err := FindNotesLinkedFrom(t.Context(), nil, &FindNotesLinkedFromInput{Path: path})
	require.NoError(t, err)
	require.Len(t, output.Notes, 1)
	assert.Equal(t, "Note 2", output.Notes[0].Frontmatter.Title)
}

func TestFindNotesLinkedFromNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := FindNotesLinkedFrom(t.Context(), nil, &FindNotesLinkedFromInput{Path: "note.md"})
	assert.Error(t, err)
}

func TestFindNotesLinkedFromNotFound(t *testing.T) {
	withVault(t)

	_, _, err := FindNotesLinkedFrom(t.Context(), nil, &FindNotesLinkedFromInput{Path: "does-not-exist.md"})
	assert.Error(t, err)
}
