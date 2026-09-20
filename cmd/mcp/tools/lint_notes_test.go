package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLintNotes(t *testing.T) {
	tmp := withVault(t)

	// A permanent note with no links is flagged as an orphan.
	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\ntype: permanent\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	_, output, err := LintNotes(t.Context(), nil, &LintNotesInput{Path: "."})
	require.NoError(t, err)
	assert.NotEmpty(t, output.Errors)
}

func TestLintNotesDefaultsPath(t *testing.T) {
	tmp := withVault(t)

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	_, output, err := LintNotes(t.Context(), nil, &LintNotesInput{})
	require.NoError(t, err)
	assert.Empty(t, output.Errors)
}

func TestLintNotesNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := LintNotes(t.Context(), nil, &LintNotesInput{Path: "."})
	assert.Error(t, err)
}
