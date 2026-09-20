package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegexSearchNotes(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody about hiking"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Note 2\n---\nBody about cooking"), 0o644)
	require.NoError(t, err)

	_, output, err := RegexSearchNotes(t.Context(), nil, &RegexSearchNotesInput{Path: ".", Expression: "hiking"})
	require.NoError(t, err)
	require.Len(t, output.Notes, 1)
	assert.Equal(t, "Note 1", output.Notes[0].Frontmatter.Title)
}

func TestRegexSearchNotesNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := RegexSearchNotes(t.Context(), nil, &RegexSearchNotesInput{Path: ".", Expression: "hiking"})
	assert.Error(t, err)
}

func TestRegexSearchNotesInvalidExpression(t *testing.T) {
	withVault(t)

	_, _, err := RegexSearchNotes(t.Context(), nil, &RegexSearchNotesInput{Path: ".", Expression: "("})
	assert.Error(t, err)
}
