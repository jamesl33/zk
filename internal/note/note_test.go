package note

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkNewNote(b *testing.B) {
	for b.Loop() {
		_, err := New("testdata/20060102150405.md")
		require.NoError(b, err)
	}
}
