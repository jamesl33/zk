package note

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSummarizeWrapNoEnv(t *testing.T) {
	require.NoError(t, os.Unsetenv("FZF_PREVIEW_COLUMNS"))

	var s Summarize

	assert.Equal(t, "some content", s.wrap("some content"))
}

func TestSummarizeWrapZeroColumns(t *testing.T) {
	t.Setenv("FZF_PREVIEW_COLUMNS", "0")

	var s Summarize

	assert.Equal(t, "some content", s.wrap("some content"))
}

func TestSummarizeWrapInvalidColumns(t *testing.T) {
	t.Setenv("FZF_PREVIEW_COLUMNS", "not-a-number")

	var s Summarize

	assert.Equal(t, "some content", s.wrap("some content"))
}

func TestSummarizeWrapColumns(t *testing.T) {
	t.Setenv("FZF_PREVIEW_COLUMNS", "10")

	var s Summarize

	assert.Equal(t, "some\ncontent\nhere", s.wrap("some content here"))
}
