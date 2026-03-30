package notes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkedFrom(t *testing.T) {
	tmp := t.TempDir()

	// We must change the working directory because LinkedFrom hardcodes "."
	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(cwd)

	// Create notes
	// note1.md has links to 20060102150405 and 20060102150406
	err = os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\n[[20060102150405]] and [[20060102150406|Other Note]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150406.md"), []byte("---\ntitle: Note 3\n---\nBody 3"), 0o644)
	require.NoError(t, err)

	n, err := note.New("20060102150404.md")
	require.NoError(t, err)

	var actual []*note.Note

	err = LinkedFrom(t.Context(), n, func(n *note.Note) {
		actual = append(actual, n)
	})

	require.NoError(t, err)
	assert.Len(t, actual, 2)

	titles := []string{actual[0].Frontmatter.Title, actual[1].Frontmatter.Title}
	assert.Contains(t, titles, "Note 2")
	assert.Contains(t, titles, "Note 3")
}

func TestLinkedFromNoLinks(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(cwd)

	err = os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\nNo links here"), 0o644)
	require.NoError(t, err)

	n, err := note.New("20060102150404.md")
	require.NoError(t, err)

	var actual []*note.Note

	err = LinkedFrom(t.Context(), n, func(n *note.Note) {
		actual = append(actual, n)
	})

	require.NoError(t, err)
	assert.Empty(t, actual)
}

func TestLinkedFromError(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(cwd)

	err = os.WriteFile(
		filepath.Join(tmp, "20060102150404.md"),
		[]byte("---\ntitle: Note 1\n---\n[[20060102150405]]"),
		0o644,
	)
	require.NoError(t, err)

	n, err := note.New("20060102150404.md")
	require.NoError(t, err)

	// Delete the file to cause an error when LinkedFrom calls Links() -> GetBody()
	require.NoError(t, os.Remove(filepath.Join(tmp, "20060102150404.md")))

	err = LinkedFrom(t.Context(), n, func(n *note.Note) {})
	assert.Error(t, err)
}

func TestLinkedTo(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(cwd)

	// note2.md and note3.md link to note1.md
	err = os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\n[[20060102150404]]"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150406.md"), []byte("---\ntitle: Note 3\n---\n[[20060102150404|Custom Text]]"), 0o644)
	require.NoError(t, err)

	n, err := note.New("20060102150404.md")
	require.NoError(t, err)

	var actual []*note.Note

	err = LinkedTo(t.Context(), n, func(n *note.Note) {
		actual = append(actual, n)
	})

	require.NoError(t, err)
	assert.Len(t, actual, 2)

	titles := []string{actual[0].Frontmatter.Title, actual[1].Frontmatter.Title}
	assert.Contains(t, titles, "Note 2")
	assert.Contains(t, titles, "Note 3")
}

func TestLinkedToNoLinks(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(cwd)

	err = os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\nNo links here"), 0o644)
	require.NoError(t, err)

	n, err := note.New("20060102150404.md")
	require.NoError(t, err)

	var actual []*note.Note

	err = LinkedTo(t.Context(), n, func(n *note.Note) {
		actual = append(actual, n)
	})

	require.NoError(t, err)
	assert.Empty(t, actual)
}

func TestLinkedToError(t *testing.T) {
	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	defer os.Chdir(cwd)

	err = os.WriteFile(filepath.Join(tmp, "20060102150404.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "20060102150405.md"), []byte("---\ntitle: Note 2\n---\n[[20060102150404]]"), 0o644)
	require.NoError(t, err)

	n, err := note.New("20060102150404.md")
	require.NoError(t, err)

	// Cause an error by making a file unreadable
	require.NoError(t, os.Chmod(filepath.Join(tmp, "20060102150405.md"), 0o000))

	err = LinkedTo(t.Context(), n, func(n *note.Note) {})
	assert.Error(t, err)
}
