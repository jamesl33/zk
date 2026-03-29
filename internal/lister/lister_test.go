package lister

import (
	"context"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/slices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const content = `---
title: "Note %d"
tags: ["tag%d"]
---

%s`

func TestNewLister(t *testing.T) {
	m := matcher.Any()

	l, err := NewLister(WithPath("/tmp"), WithMatcher(m))
	require.NoError(t, err)

	assert.Equal(t, "/tmp", l.options.path)
	assert.NotNil(t, l.options.matcher)
}

func TestListerMany(t *testing.T) {
	seed := []struct {
		name    string
		content string
	}{
		{
			name:    "note1.md",
			content: "---\ntitle: Note 1\n---\nBody 1",
		},
		{
			name:    "note2.md",
			content: "---\ntitle: Note 2\n---\nBody 2",
		},
		{
			name:    ".hidden.md",
			content: "---\ntitle: Hidden\n---\nBody",
		},
		{
			name:    "GEMINI.md",
			content: "---\ntitle: Gemini\n---\nBody",
		},
		{
			name:    "not-a-note.txt",
			content: "just text",
		},
		{
			name:    "invalid-frontmatter.md",
			content: "---\ninvalid\n---\nBody",
		},
	}

	tmp := t.TempDir()

	for _, n := range seed {
		err := os.WriteFile(filepath.Join(tmp, n.name), []byte(n.content), 0o644)
		require.NoError(t, err)
	}

	l, err := NewLister(WithPath(tmp))
	require.NoError(t, err)

	actual, err := slices.Collect2[[]*note.Note](l.Many(t.Context()))
	require.NoError(t, err)

	expected := []*note.Note{
		{
			Path:        filepath.Join(tmp, "note1.md"),
			Frontmatter: note.Frontmatter{Title: "Note 1"},
		},
		{
			Path:        filepath.Join(tmp, "note2.md"),
			Frontmatter: note.Frontmatter{Title: "Note 2"},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestListerOne(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	l, err := NewLister(WithPath(tmp))
	require.NoError(t, err)

	actual, err := l.One(t.Context())
	require.NoError(t, err)

	expected := &note.Note{
		Path:        filepath.Join(tmp, "note1.md"),
		Frontmatter: note.Frontmatter{Title: "Note 1"},
	}

	assert.Equal(t, expected, actual)
}

func TestListerManyWithMatcher(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note2.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	m := func(n *note.Note) (bool, error) {
		return n.Frontmatter.Title == "Note 2", nil
	}

	l, err := NewLister(WithPath(tmp), WithMatcher(m))
	require.NoError(t, err)

	actual, err := slices.Collect2[[]*note.Note](l.Many(t.Context()))
	require.NoError(t, err)

	expected := []*note.Note{
		{
			Path:        filepath.Join(tmp, "note2.md"),
			Frontmatter: note.Frontmatter{Title: "Note 2"},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestListerManyWithRecursion(t *testing.T) {
	tmp := t.TempDir()

	err := os.MkdirAll(filepath.Join(tmp, "subdir", "subsubdir"), 0o755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "note1.md"), []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "subdir", "note2.md"), []byte("---\ntitle: Note 2\n---\nBody 2"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmp, "subdir", "subsubdir", "note3.md"), []byte("---\ntitle: Note 3\n---\nBody 3"), 0o644)
	require.NoError(t, err)

	l, err := NewLister(WithPath(tmp))
	require.NoError(t, err)

	actual, err := slices.Collect2[[]*note.Note](l.Many(t.Context()))
	require.NoError(t, err)

	expected := []*note.Note{
		{
			Path:        filepath.Join(tmp, "note1.md"),
			Frontmatter: note.Frontmatter{Title: "Note 1"},
		},
		{
			Path:        filepath.Join(tmp, "subdir", "note2.md"),
			Frontmatter: note.Frontmatter{Title: "Note 2"},
		},
		{
			Path:        filepath.Join(tmp, "subdir", "subsubdir", "note3.md"),
			Frontmatter: note.Frontmatter{Title: "Note 3"},
		},
	}

	assert.Equal(t, expected, actual)
}

func TestListerManyWithContextCancellation(t *testing.T) {
	tmp := t.TempDir()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	l, err := NewLister(WithPath(tmp))
	require.NoError(t, err)

	next, stop := iter.Pull2(l.Many(ctx))
	defer stop()

	_, err, ok := next()
	assert.True(t, ok)
	assert.ErrorIs(t, err, context.Canceled)
}

func BenchmarkListerVariyingNotes(b *testing.B) {
	for _, notes := range []int{10, 100, 1000, 10000, 100000} {
		b.Run(fmt.Sprintf("%d-notes", notes), func(b *testing.B) {
			tmp := b.TempDir()

			for i := 0; i < notes; i++ {
				var (
					path    = filepath.Join(tmp, fmt.Sprintf("%d.md", i))
					content = fmt.Sprintf(content, i, i%10, strconv.Itoa(i))
				)

				err := os.WriteFile(path, []byte(content), 0o644)
				require.NoError(b, err)
			}

			l, err := NewLister(WithPath(tmp))
			require.NoError(b, err)

			b.ResetTimer()

			for b.Loop() {
				iter, stop := iter.Pull2(l.Many(b.Context()))
				defer stop()

				count := 0

				for {
					n, err, ok := iter()
					if !ok {
						break
					}

					require.NoError(b, err)
					require.NotNil(b, n)
					count++
				}

				require.Equal(b, notes, count)
			}
		})
	}
}

func BenchmarkListerVaryingSize(b *testing.B) {
	const notes = 250

	for _, size := range []int{1024, 1024 * 100, 1024 * 1024, 10 * 1024 * 1024, 20 * 1024 * 1024} {
		b.Run(fmt.Sprintf("%d-bytes", size), func(b *testing.B) {
			var (
				tmp  = b.TempDir()
				body = strings.Repeat("a", size)
			)

			for i := 0; i < notes; i++ {
				var (
					path    = filepath.Join(tmp, fmt.Sprintf("%d.md", i))
					content = fmt.Sprintf(content, i, i%10, body)
				)

				err := os.WriteFile(path, []byte(content), 0o644)
				require.NoError(b, err)
			}

			l, err := NewLister(WithPath(tmp))
			require.NoError(b, err)

			b.ResetTimer()

			for b.Loop() {
				iter, stop := iter.Pull2(l.Many(b.Context()))
				defer stop()

				count := 0

				for {
					n, err, ok := iter()
					if !ok {
						break
					}

					require.NoError(b, err)
					require.NotNil(b, n)

					count++
				}

				require.Equal(b, notes, count)
			}
		})
	}
}

func BenchmarkListerVaryingSizeWithABodyMatcher(b *testing.B) {
	const notes = 250

	for _, size := range []int{1024, 1024 * 100, 1024 * 1024, 10 * 1024 * 1024, 20 * 1024 * 1024} {
		b.Run(fmt.Sprintf("%d-bytes", size), func(b *testing.B) {
			var (
				tmp  = b.TempDir()
				body = strings.Repeat("a", size)
			)

			for i := 0; i < notes; i++ {
				var (
					path    = filepath.Join(tmp, fmt.Sprintf("%d.md", i))
					content = fmt.Sprintf(content, i, i%10, body)
				)

				err := os.WriteFile(path, []byte(content), 0o644)
				require.NoError(b, err)
			}

			matcher, err := matcher.Body("", "", ".*")
			require.NoError(b, err)

			l, err := NewLister(WithPath(tmp), WithMatcher(matcher))
			require.NoError(b, err)

			b.ResetTimer()

			for b.Loop() {
				iter, stop := iter.Pull2(l.Many(b.Context()))
				defer stop()

				count := 0

				for {
					n, err, ok := iter()
					if !ok {
						break
					}

					require.NoError(b, err)
					require.NotNil(b, n)

					count++
				}

				require.Equal(b, notes, count)
			}
		})
	}
}
