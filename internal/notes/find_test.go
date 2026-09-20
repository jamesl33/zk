package notes

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeIndex is a minimal in-memory 'index' used to test Finder without a real vector.DB/AI client.
type fakeIndex struct {
	upserted []string
	results  []*note.Note
	findErr  error
	closed   bool
}

func (f *fakeIndex) Upsert(_ context.Context, n *note.Note) error {
	f.upserted = append(f.upserted, n.Name())
	return nil
}

func (f *fakeIndex) Find(_ context.Context, _ *note.Note) ([]*note.Note, error) {
	return f.results, f.findErr
}

func (f *fakeIndex) Close() error {
	f.closed = true
	return nil
}

// newTestNote writes and opens a note in the given directory.
func newTestNote(t *testing.T, dir, name, content string) *note.Note {
	t.Helper()

	path := filepath.Join(dir, name+".md")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	n, err := note.New(path)
	require.NoError(t, err)

	return n
}

// chdir changes into the given directory for the duration of the test.
func chdir(t *testing.T, dir string) {
	t.Helper()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })
}

func TestFinderFindPopulatesAndReturnsResults(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	newTestNote(t, tmp, "20240101000001", "---\ntitle: Note 1\n---\nBody 1")
	target := newTestNote(t, tmp, "20240101000002", "---\ntitle: Note 2\n---\nBody 2")

	idx := &fakeIndex{results: []*note.Note{target}}
	finder := &Finder{db: idx}

	var found []*note.Note

	err := finder.Find(t.Context(), target, func(n *note.Note) error {
		found = append(found, n)
		return nil
	})
	require.NoError(t, err)
	require.Len(t, found, 1)
	assert.Equal(t, target.Name(), found[0].Name())
	assert.Contains(t, idx.upserted, "20240101000001")
}

func TestFinderFindPropagatesError(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	n := newTestNote(t, tmp, "20240101000001", "---\ntitle: Note 1\n---\nBody 1")

	idx := &fakeIndex{findErr: assert.AnError}
	finder := &Finder{db: idx}

	err := finder.Find(t.Context(), n, func(n *note.Note) error { return nil })
	assert.ErrorIs(t, err, assert.AnError)
}

func TestFinderClose(t *testing.T) {
	idx := &fakeIndex{}
	finder := &Finder{db: idx}

	require.NoError(t, finder.Close())
	assert.True(t, idx.closed)
}
