package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListNotes(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	_, output, err := ListNotes(t.Context(), nil, &ListNotesInput{Path: "."})
	require.NoError(t, err)
	assert.Len(t, output.Notes, 2)
}

func TestListNotesNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := ListNotes(t.Context(), nil, &ListNotesInput{Path: "."})
	assert.Error(t, err)
}
