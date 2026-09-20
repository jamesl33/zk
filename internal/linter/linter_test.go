package linter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinter(t *testing.T) {
	l := NewLinter()
	assert.NotNil(t, l)
}

func TestLinterLintNoErrors(t *testing.T) {
	var (
		tmp = t.TempDir()
		// note1 links to note2
		note1 = filepath.Join(tmp, "20240101000001.md")
		// note2 has no links
		note2 = filepath.Join(tmp, "20240101000002.md")
	)

	err := os.WriteFile(note1, []byte("---\ntitle: Note 1\ndate: \"2024-01-01\"\n---\n[[20240101000002]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(note2, []byte("---\ntitle: Note 2\ndate: \"2024-01-01\"\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestLinterLintBrokenLink(t *testing.T) {
	var (
		tmp = t.TempDir()
		// note1 links to missing note
		note1 = filepath.Join(tmp, "20240101000001.md")
	)

	err := os.WriteFile(note1, []byte("---\ntitle: Note 1\ndate: \"2024-01-01\"\n---\n[[20240101000002]]"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	require.Len(t, errors, 1)
	assert.Equal(t, note1, errors[0].Path)
	assert.Equal(t, fmt.Sprintf("Link %q is broken (linkcheck)", "20240101000002"), errors[0].Message)
	assert.Equal(t, 5, errors[0].Line)
	assert.Equal(t, 1, errors[0].Column)
}

func TestLinterLintBrokenLinkMidLine(t *testing.T) {
	var (
		tmp   = t.TempDir()
		note1 = filepath.Join(tmp, "20240101000001.md")
	)

	err := os.WriteFile(note1, []byte("---\ntitle: Note 1\ndate: \"2024-01-01\"\n---\nSee [[20240101000002]] above"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	require.Len(t, errors, 1)
	assert.Equal(t, 5, errors[0].Line)
	assert.Equal(t, 5, errors[0].Column)
}

func TestLinterLintMultipleBrokenLinks(t *testing.T) {
	var (
		tmp = t.TempDir()
		// note1 links to two missing notes
		note1 = filepath.Join(tmp, "20240101000001.md")
	)

	err := os.WriteFile(note1, []byte("---\ntitle: Note 1\ndate: \"2024-01-01\"\n---\n[[20240101000002]] [[20240101000003]]"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	require.Len(t, errors, 2)
}

func TestLinterLintLinkToSelf(t *testing.T) {
	var (
		tmp   = t.TempDir()
		note1 = filepath.Join(tmp, "20240101000001.md")
	)

	err := os.WriteFile(note1, []byte("---\ntitle: Note 1\ndate: \"2024-01-01\"\n---\n[[20240101000001]]"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestLinterLintDuplicateIdentifier(t *testing.T) {
	var (
		tmp   = t.TempDir()
		note1 = filepath.Join(tmp, "20240101000001.md")
		note2 = filepath.Join(tmp, "subdir", "20240101000001.md")
	)

	err := os.MkdirAll(filepath.Join(tmp, "subdir"), 0o755)
	require.NoError(t, err)

	err = os.WriteFile(note1, []byte("---\ntitle: Note 1\ndate: \"2024-01-01\"\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(note2, []byte("---\ntitle: Note 2\ndate: \"2024-01-01\"\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	require.Len(t, errors, 2)

	paths := []string{errors[0].Path, errors[1].Path}
	assert.ElementsMatch(t, []string{note1, note2}, paths)

	for _, e := range errors {
		assert.Contains(t, e.Message, "duplicate-id")
	}
}

func TestLinterLintNestedDirectories(t *testing.T) {
	var (
		tmp   = t.TempDir()
		note1 = filepath.Join(tmp, "20240101000001.md")
		note2 = filepath.Join(tmp, "subdir", "20240101000002.md")
	)

	err := os.MkdirAll(filepath.Join(tmp, "subdir"), 0o755)
	require.NoError(t, err)

	err = os.WriteFile(note1, []byte("---\ntitle: Note 1\ndate: \"2024-01-01\"\n---\n[[20240101000002]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(note2, []byte("---\ntitle: Note 2\ndate: \"2024-01-01\"\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestLinterLintEmptyDirectory(t *testing.T) {
	tmp := t.TempDir()

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestLinterLintNoLinks(t *testing.T) {
	var (
		tmp   = t.TempDir()
		note1 = filepath.Join(tmp, "20240101000001.md")
	)

	err := os.WriteFile(note1, []byte("---\ntitle: Note 1\ndate: \"2024-01-01\"\n---\nNo links here"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestLinterLintWithLinkText(t *testing.T) {
	var (
		tmp   = t.TempDir()
		note1 = filepath.Join(tmp, "20240101000001.md")
		note2 = filepath.Join(tmp, "20240101000002.md")
	)

	err := os.WriteFile(note1, []byte("---\ntitle: Note 1\ndate: \"2024-01-01\"\n---\n[[20240101000002|Custom Text]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(note2, []byte("---\ntitle: Note 2\ndate: \"2024-01-01\"\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestLinterLintOrphanPermanentNote(t *testing.T) {
	var (
		tmp   = t.TempDir()
		note1 = filepath.Join(tmp, "20240101000001.md")
	)

	err := os.WriteFile(note1, []byte("---\ntype: permanent\ntitle: Note 1\ndate: \"2024-01-01\"\n---\nNo links here"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	require.Len(t, errors, 1)
	assert.Equal(t, note1, errors[0].Path)
	assert.Equal(t, "Permanent note has no links (orphan-note)", errors[0].Message)
}

func TestLinterLintPermanentNoteWithLink(t *testing.T) {
	var (
		tmp   = t.TempDir()
		note1 = filepath.Join(tmp, "20240101000001.md")
		note2 = filepath.Join(tmp, "20240101000002.md")
	)

	err := os.WriteFile(note1, []byte("---\ntype: permanent\ntitle: Note 1\ndate: \"2024-01-01\"\n---\n[[20240101000002]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(note2, []byte("---\ntitle: Note 2\ndate: \"2024-01-01\"\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestLinterLintNonPermanentNoteWithNoLinks(t *testing.T) {
	var (
		tmp   = t.TempDir()
		note1 = filepath.Join(tmp, "20240101000001.md")
	)

	err := os.WriteFile(note1, []byte("---\ntype: fleeting\ntitle: Note 1\ndate: \"2024-01-01\"\n---\nNo links here"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestLinterLintExcludesArchivesByDefault(t *testing.T) {
	var (
		tmp     = t.TempDir()
		archive = filepath.Join(tmp, archiveDir)
		note1   = filepath.Join(archive, "20240101000001.md")
	)

	require.NoError(t, os.MkdirAll(archive, 0o755))

	err := os.WriteFile(note1, []byte("---\ntype: permanent\ntitle: Note 1\ndate: \"2024-01-01\"\n---\nNo links here"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp)
	require.NoError(t, err)
	assert.Empty(t, errors)
}

func TestLinterLintIncludesArchivesWithOption(t *testing.T) {
	var (
		tmp     = t.TempDir()
		archive = filepath.Join(tmp, archiveDir)
		note1   = filepath.Join(archive, "20240101000001.md")
	)

	require.NoError(t, os.MkdirAll(archive, 0o755))

	err := os.WriteFile(note1, []byte("---\ntype: permanent\ntitle: Note 1\ndate: \"2024-01-01\"\n---\nNo links here"), 0o644)
	require.NoError(t, err)

	l := NewLinter()

	errors, err := l.Lint(t.Context(), tmp, WithArchives())
	require.NoError(t, err)
	require.Len(t, errors, 1)
	assert.Equal(t, note1, errors[0].Path)
	assert.Equal(t, "Permanent note has no links (orphan-note)", errors[0].Message)
}

func TestLinterLintInvalidPath(t *testing.T) {
	l := NewLinter()

	_, err := l.Lint(context.Background(), "/non-existent-path-123")
	assert.Error(t, err)
}

func TestLinterLintContextCanceled(t *testing.T) {
	tmp := t.TempDir()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	l := NewLinter()

	_, err := l.Lint(ctx, tmp)
	assert.Error(t, err)
}
