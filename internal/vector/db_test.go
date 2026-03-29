package vector

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

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
		checksum  uint32
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

	// Find lists within the current working directory, so we must be within that directory
	cwd, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(tmp)
	require.NoError(t, err)
	defer os.Chdir(cwd)

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
