package cache

import (
	"hash/crc32"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/mattn/go-sqlite3"
)

func TestNew(t *testing.T) {
	var (
		tmp  = t.TempDir()
		path = filepath.Join(tmp, "cache.db")
	)

	c, err := New[string](t.Context(), path, "test_table")
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c.db)
	assert.Equal(t, "test_table", c.table)
}

func TestNewError(t *testing.T) {
	// sqlite doesn't like table names starting with numbers without quoting.
	c, err := New[string](t.Context(), ":memory:", "1invalid")
	assert.Error(t, err)
	assert.Nil(t, c)
}

func TestGet(t *testing.T) {
	var (
		tmp  = t.TempDir()
		path = filepath.Join(tmp, "cache.db")
	)

	c, err := New[string](t.Context(), path, "test_table")
	require.NoError(t, err)

	err = c.Set(t.Context(), "prompt", "result")
	require.NoError(t, err)

	actual, err := c.Get(t.Context(), "prompt")
	require.NoError(t, err)
	require.NotNil(t, actual)
	assert.Equal(t, "result", *actual)
}

func TestGetNotFound(t *testing.T) {
	var (
		tmp  = t.TempDir()
		path = filepath.Join(tmp, "cache.db")
	)

	c, err := New[string](t.Context(), path, "test_table")
	require.NoError(t, err)

	actual, err := c.Get(t.Context(), "nonexistent")
	require.NoError(t, err)
	assert.Nil(t, actual)
}

func TestGetError(t *testing.T) {
	var (
		tmp  = t.TempDir()
		path = filepath.Join(tmp, "cache.db")
	)

	c, err := New[string](t.Context(), path, "test_table")
	require.NoError(t, err)

	// drop the table to cause an error on Get
	_, err = c.db.Exec("DROP TABLE test_table")
	require.NoError(t, err)

	actual, err := c.Get(t.Context(), "prompt")
	assert.Error(t, err)
	assert.Nil(t, actual)
}

func TestSet(t *testing.T) {
	var (
		tmp  = t.TempDir()
		path = filepath.Join(tmp, "cache.db")
	)

	c, err := New[[]byte](t.Context(), path, "test_table")
	require.NoError(t, err)

	expected := []byte("result")

	err = c.Set(t.Context(), "prompt", expected)
	assert.NoError(t, err)

	var (
		actual []byte
		hasher = crc32.NewIEEE()
	)

	_, err = io.Copy(hasher, strings.NewReader("prompt"))
	require.NoError(t, err)

	err = c.db.QueryRow("SELECT value FROM test_table WHERE key = ?", hasher.Sum32()).Scan(&actual)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSetError(t *testing.T) {
	var (
		tmp  = t.TempDir()
		path = filepath.Join(tmp, "cache.db")
	)

	c, err := New[string](t.Context(), path, "test_table")
	require.NoError(t, err)

	// drop the table to cause an error on Set
	_, err = c.db.Exec("DROP TABLE test_table")
	require.NoError(t, err)

	err = c.Set(t.Context(), "prompt", "result")
	assert.Error(t, err)
}
