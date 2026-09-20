package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindNotesLinkedTo(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\n[[20060102150404]]"), 0o644)
	require.NoError(t, err)

	path := filepath.Join(tmp, "20060102150404.md")

	_, output, err := FindNotesLinkedTo(t.Context(), nil, &FindNotesLinkedToInput{Path: path})
	require.NoError(t, err)
	require.Len(t, output.Notes, 1)
	assert.Equal(t, "Note 2", output.Notes[0].Frontmatter.Title)
}

func TestFindNotesLinkedToNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := FindNotesLinkedTo(t.Context(), nil, &FindNotesLinkedToInput{Path: "note.md"})
	assert.Error(t, err)
}

func TestFindNotesLinkedToNotFound(t *testing.T) {
	withVault(t)

	_, _, err := FindNotesLinkedTo(t.Context(), nil, &FindNotesLinkedToInput{Path: "does-not-exist.md"})
	assert.Error(t, err)
}
