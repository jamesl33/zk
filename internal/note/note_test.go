package note

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNote(t *testing.T) {
	_, err := New("testdata/20060102150405.md")
	require.NoError(t, err)
}

func TestNewNoteNotANote(t *testing.T) {
	_, err := New("testdata/20060102150406.md")
	require.ErrorIs(t, err, ErrNotNote)
}

func BenchmarkNewNote(b *testing.B) {
	for b.Loop() {
		_, err := New("testdata/20060102150405.md")
		require.NoError(b, err)
	}
}
