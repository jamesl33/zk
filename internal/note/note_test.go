package note

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	n, err := New("testdata/20060102150405.md")
	require.NoError(t, err)

	assert.Equal(t, "testdata/20060102150405.md", n.Path)
	assert.Equal(t, Type("permanent"), n.Frontmatter.Type)
	assert.Equal(t, "Lorem Ipsum", n.Frontmatter.Title)
}

func TestNewNotNote(t *testing.T) {
	_, err := New("testdata/20060102150406.md")
	assert.ErrorIs(t, err, ErrNotNote)
}

func TestNewFileNotFound(t *testing.T) {
	_, err := New("testdata/does-not-exist.md")
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestNoteName(t *testing.T) {
	n := Note{Path: "/path/to/my-note.md"}
	assert.Equal(t, "my-note", n.Name())
}

func TestNoteGetBody(t *testing.T) {
	n, err := New("testdata/20060102150405.md")
	require.NoError(t, err)

	body, err := n.GetBody()
	require.NoError(t, err)
	assert.Contains(t, body, "Lorem ipsum dolor sit amet")
}

func TestNoteGetBodyCached(t *testing.T) {
	n := Note{
		body: ptr.To("Cached body"),
	}

	body, err := n.GetBody()
	require.NoError(t, err)
	assert.Equal(t, "Cached body", body)
}

func TestNoteSetBody(t *testing.T) {
	var n Note

	n.SetBody("New body")

	assert.Equal(t, "New body", *n.body)
}

func TestNoteLinks(t *testing.T) {
	tmp := t.TempDir()

	path := filepath.Join(tmp, "note.md")

	// Links must be 14 digits
	content := "---\ntype: permanent\ntitle: Note\n---\n[[20251129123456]] and [[20251129123457]]"

	err := os.WriteFile(path, []byte(content), 0o644)
	require.NoError(t, err)

	n, err := New(path)
	require.NoError(t, err)

	links, err := n.Links()
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"20251129123456", "20251129123457"}, links)
}

func TestNoteChecksum(t *testing.T) {
	n, err := New("testdata/20060102150405.md")
	require.NoError(t, err)

	first, err := n.Checksum()
	require.NoError(t, err)
	assert.NotZero(t, first)

	second, err := n.Checksum()
	require.NoError(t, err)
	assert.NotZero(t, second)

	require.Equal(t, first, second)
}

func TestNoteWrite(t *testing.T) {
	var (
		tmp  = t.TempDir()
		path = filepath.Join(tmp, "note.md")
	)

	n := Note{
		Path:        path,
		Frontmatter: Frontmatter{Type: TypePermanent, Title: "Test Note"},
	}

	n.SetBody("Test body")

	err := n.Write()
	require.NoError(t, err)

	// Verify by reading back
	n2, err := New(path)
	require.NoError(t, err)
	assert.Equal(t, "Test Note", n2.Frontmatter.Title)

	body, err := n2.GetBody()
	require.NoError(t, err)
	assert.Equal(t, "Test body", body)
}

func TestNoteText(t *testing.T) {
	n, err := New("testdata/20060102150405.md")
	require.NoError(t, err)

	text, err := n.Text()
	require.NoError(t, err)
	assert.Contains(t, text, "---")
	assert.Contains(t, text, "title: Lorem Ipsum")
	assert.Contains(t, text, "Lorem ipsum dolor sit amet")
}

func TestNoteString(t *testing.T) {
	n := Note{
		Path:        "notes/permanent/note.md",
		Frontmatter: Frontmatter{Title: "My Note", Tags: []string{"tag1", "tag2"}},
	}

	s := n.String()

	// String() output contains ANSI color codes because init() forces NoColor = false.
	assert.Contains(t, s, "notes/permanent")
	assert.Contains(t, s, "My Note")
	assert.Contains(t, s, "tag1,tag2")
}

func TestNoteGetBodyFileNotFound(t *testing.T) {
	n := Note{
		Path: "does-not-exist.md",
	}

	_, err := n.GetBody()
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestNoteWriteError(t *testing.T) {
	n := Note{
		Path: "/non-existent-dir/note.md",
	}

	err := n.Write()
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestNoteWriteTo(t *testing.T) {
	var (
		tmp  = t.TempDir()
		path = filepath.Join(tmp, "note.md")
	)

	n := Note{
		Path:        path,
		Frontmatter: Frontmatter{Type: TypePermanent, Title: "Test Note"},
	}

	n.SetBody("Test body")

	var b bytes.Buffer

	_, err := n.WriteTo(&b)
	require.NoError(t, err)

	content := b.String()
	assert.Contains(t, content, "---\n")
	assert.Contains(t, content, "title: Test Note")
	assert.Contains(t, content, "Test body")
}

func TestNoteEdit(t *testing.T) {
	var (
		tmp     = t.TempDir()
		path    = filepath.Join(tmp, "note.md")
		content = "---\ntype: permanent\ntitle: Old Title\n---\nOld body"
	)

	err := os.WriteFile(path, []byte(content), 0o644)
	require.NoError(t, err)

	n, err := New(path)
	require.NoError(t, err)

	// Create a script to use as the editor
	editor := filepath.Join(tmp, "editor.sh")

	err = os.WriteFile(editor, []byte("#!/bin/sh\nsed -i 's/Old Title/New Title/' \"$1\""), 0o755)
	require.NoError(t, err)

	os.Setenv("EDITOR", editor)
	defer os.Unsetenv("EDITOR")

	err = n.Edit(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "New Title", n.Frontmatter.Title)
}

func TestNoteEditNoEditor(t *testing.T) {
	// Ensure EDITOR is not set
	os.Setenv("EDITOR", "")
	defer os.Unsetenv("EDITOR")

	var n Note

	err := n.Edit(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no editor set")
}

func BenchmarkNewNote(b *testing.B) {
	for b.Loop() {
		_, err := New("testdata/20060102150405.md")
		require.NoError(b, err)
	}
}

func BenchmarkNoteString(b *testing.B) {
	n := Note{
		Path:        "notes/permanent/note.md",
		Frontmatter: Frontmatter{Title: "My Note", Tags: []string{"tag1", "tag2"}},
	}

	b.ResetTimer()

	for b.Loop() {
		_ = n.String()
	}
}
