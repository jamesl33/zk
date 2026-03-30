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
