package vector

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	mock_ai "github.com/jamesl33/zk/internal/ai/mocks"
	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// newNote returns a new note in the given directory, with the provided name/content.
func newNote(t *testing.T, tmp, name, content string) *note.Note {
	path := filepath.Join(tmp, name+".md")

	err := os.WriteFile(path, []byte(content), 0o644)
	require.NoError(t, err)

	n, err := note.New(path)
	require.NoError(t, err)

	return n
}

func TestNew(t *testing.T) {
	tmp := t.TempDir()

	db, err := New(t.Context(), filepath.Join(tmp, "zk.sqlite3"))
	require.NoError(t, err)
	defer db.Close()

	var name string

	err = db.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='notes'").Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "notes", name)
}

func TestNewClientFailure(t *testing.T) {
	tmp := t.TempDir()

	// Passing a directory as the database path causes cache/client creation to fail.
	_, err := New(t.Context(), tmp)
	assert.Error(t, err)
}

func TestDBInit(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	vdb := &DB{
		db: db,
	}

	err = vdb.init(t.Context())
	require.NoError(t, err)

	var name string

	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='notes'").Scan(&name)
	require.NoError(t, err)
	assert.Equal(t, "notes", name)
}

func TestDBInitFailure(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	err = db.Close() // Close it to force failure
	require.NoError(t, err)

	vdb := &DB{
		db: db,
	}

	err = vdb.init(t.Context())
	assert.Error(t, err)
}

func TestDBClose(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	vdb := &DB{
		db: db,
	}

	err = vdb.Close()
	assert.NoError(t, err)
}

func TestDBUpsertSuccess(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	vdb := &DB{
		client: mclient,
		db:     db,
	}

	require.NoError(t, vdb.init(t.Context()))

	n := newNote(t, tmp, "note1", "---\ntitle: Note 1\n---\nBody 1")

	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return([]float32{1.0, 2.0}, nil)

	err = vdb.Upsert(t.Context(), n)
	require.NoError(t, err)

	var (
		name      string
		checksum  []byte
		embedding []byte
	)

	err = db.QueryRow("SELECT name, checksum, embedding FROM notes WHERE name = ?", n.Name()).Scan(&name, &checksum, &embedding)
	require.NoError(t, err)

	assert.Equal(t, n.Name(), name)

	cs, err := n.Checksum()
	require.NoError(t, err)
	assert.Equal(t, cs, checksum)
}

func TestDBUpsertSkip(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	vdb := &DB{
		client: mclient,
		db:     db,
	}

	require.NoError(t, vdb.init(t.Context()))

	n := newNote(t, tmp, "note1", "---\ntitle: Note 1\n---\nBody 1")

	// First time: insert
	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return([]float32{1.0, 2.0}, nil)

	err = vdb.Upsert(t.Context(), n)
	require.NoError(t, err)

	// Second time: skip (no Embed call expected)
	err = vdb.Upsert(t.Context(), n)
	require.NoError(t, err)
}

func TestDBUpsertMultiChunk(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	vdb := &DB{
		client: mclient,
		db:     db,
	}

	require.NoError(t, vdb.init(t.Context()))

	// Two paragraphs, each of which fits under the chunk budget alone, but not together -- forces
	// the note to be split into two chunks.
	para := strings.Repeat("filler ", 500)
	n := newNote(t, tmp, "note1", "---\ntitle: Note 1\n---\n"+para+"\n\n"+para)

	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return([]float32{1.0, 2.0}, nil).
		Times(2)

	err = vdb.Upsert(t.Context(), n)
	require.NoError(t, err)

	rows, err := db.Query("SELECT chunk FROM notes WHERE name = ? ORDER BY chunk", n.Name())
	require.NoError(t, err)
	defer rows.Close()

	var chunks []int

	for rows.Next() {
		var c int
		require.NoError(t, rows.Scan(&c))
		chunks = append(chunks, c)
	}

	require.NoError(t, rows.Err())
	assert.Equal(t, []int{0, 1}, chunks)
}

func TestDBUpsertShrinkRemovesOrphanChunks(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	vdb := &DB{
		client: mclient,
		db:     db,
	}

	require.NoError(t, vdb.init(t.Context()))

	path := filepath.Join(tmp, "note1.md")

	para := strings.Repeat("filler ", 500)
	require.NoError(t, os.WriteFile(path, []byte("---\ntitle: Note 1\n---\n"+para+"\n\n"+para), 0o644))

	big, err := note.New(path)
	require.NoError(t, err)

	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return([]float32{1.0, 2.0}, nil).
		Times(2)

	require.NoError(t, vdb.Upsert(t.Context(), big))

	var count int
	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM notes WHERE name = ?", big.Name()).Scan(&count))
	require.Equal(t, 2, count)

	// Shrink the note back down to a single chunk; the stale second chunk row must not linger.
	require.NoError(t, os.WriteFile(path, []byte("---\ntitle: Note 1\n---\nBody 1"), 0o644))

	small, err := note.New(path)
	require.NoError(t, err)

	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return([]float32{3.0, 4.0}, nil)

	require.NoError(t, vdb.Upsert(t.Context(), small))

	require.NoError(t, db.QueryRow("SELECT COUNT(*) FROM notes WHERE name = ?", big.Name()).Scan(&count))
	assert.Equal(t, 1, count)
}

func TestDBUpsertFailure(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	vdb := &DB{
		client: mclient,
		db:     db,
	}

	require.NoError(t, vdb.init(t.Context()))

	n := newNote(t, tmp, "note1", "---\ntitle: Note 1\n---\nBody 1")

	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return(nil, assert.AnError)

	err = vdb.Upsert(t.Context(), n)
	assert.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestDBFindSuccess(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	// Find lists within the vault root, so we must be within that directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmp)
	require.NoError(t, err)
	defer os.Chdir(cwd)

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	vdb := &DB{
		client: mclient,
		db:     db,
	}

	require.NoError(t, vdb.init(t.Context()))

	// Note 2: already indexed
	n2 := newNote(t, tmp, "note2", "---\ntitle: Note 2\n---\nBody 2")

	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return([]float32{1.0, 1.0}, nil)

	err = vdb.Upsert(t.Context(), n2)
	require.NoError(t, err)

	// Note 1: the one we are searching with
	n1 := newNote(t, tmp, "note1", "---\ntitle: Note 1\n---\nBody 1")

	// Similar embedding to n2
	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return([]float32{1.0, 1.1}, nil)

	results, err := vdb.Find(t.Context(), n1)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, n2.Name(), results[0].Name())
}

func TestDBFindCollapsesToClosestChunk(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	// Find lists within the vault root, so we must be within that directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmp)
	require.NoError(t, err)
	defer os.Chdir(cwd)

	require.NoError(t, os.Mkdir(filepath.Join(tmp, ".zk"), 0o755))

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	vdb := &DB{
		client: mclient,
		db:     db,
	}

	require.NoError(t, vdb.init(t.Context()))

	// note2 exists on disk, but its index rows are inserted directly rather than via Upsert, to
	// simulate a note that was chunked into two pieces: one far from the query embedding (would
	// fail the distance threshold alone), one close (should still make the note match overall).
	n2 := newNote(t, tmp, "note2", "---\ntitle: Note 2\n---\nBody 2")

	far, err := sqlite_vec.SerializeFloat32([]float32{1.0, 0.0})
	require.NoError(t, err)

	near, err := sqlite_vec.SerializeFloat32([]float32{1.0, 1.0})
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO notes VALUES (?, 0, ?, ?);`, n2.Name(), []byte("cs"), far)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO notes VALUES (?, 1, ?, ?);`, n2.Name(), []byte("cs"), near)
	require.NoError(t, err)

	n1 := newNote(t, tmp, "note1", "---\ntitle: Note 1\n---\nBody 1")

	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return([]float32{1.0, 1.1}, nil)

	results, err := vdb.Find(t.Context(), n1)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, n2.Name(), results[0].Name())
}

func TestDBFindFailure(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	vdb := &DB{
		client: mclient,
		db:     db,
	}

	require.NoError(t, vdb.init(t.Context()))

	n := newNote(t, tmp, "note1", "---\ntitle: Note 1\n---\nBody 1")

	mclient.
		EXPECT().
		Embed(gomock.Any(), gomock.Any()).
		Return(nil, assert.AnError)

	results, err := vdb.Find(t.Context(), n)
	assert.Error(t, err)
	assert.ErrorIs(t, err, assert.AnError)
	assert.Nil(t, results)
}
