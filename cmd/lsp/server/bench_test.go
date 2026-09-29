package server

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/stretchr/testify/require"
)

// benchVault creates a vault holding the given number of notes, each with a small body, and changes into it.
func benchVault(b *testing.B, notes int) {
	b.Helper()

	dir := b.TempDir()

	for i := range notes {
		var (
			name    = fmt.Sprintf("2024%010d", i)
			content = fmt.Sprintf("---\ntype: permanent\ntitle: Note %d\ntags: [a, b]\n---\n\n%s\n", i, "Some body text. [[20240000000000]]")
		)

		require.NoError(b, os.WriteFile(filepath.Join(dir, name+".md"), []byte(content), 0o644))
	}

	require.NoError(b, os.Mkdir(filepath.Join(dir, ".zk"), 0o755))

	cwd, err := os.Getwd()
	require.NoError(b, err)

	require.NoError(b, os.Chdir(dir))
	b.Cleanup(func() { require.NoError(b, os.Chdir(cwd)) })
}

// BenchmarkDiagnostics measures the cost of a single keystroke in a note containing links, which walks the vault.
func BenchmarkDiagnostics(b *testing.B) {
	for _, notes := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("notes=%d", notes), func(b *testing.B) {
			benchVault(b, notes)

			s, err := NewServer(b.Context())
			require.NoError(b, err)

			require.NoError(b, os.WriteFile("source.md", []byte("See [[20240000000001]]"), 0o644))

			uri := "file://" + filepath.Join(mustGetwd(b), "source.md")

			b.ResetTimer()

			for b.Loop() {
				_, err := s.diagnostics(uri)
				require.NoError(b, err)
			}
		})
	}
}

// BenchmarkCompletion measures the cost of a single completion request, which walks the vault.
func BenchmarkCompletion(b *testing.B) {
	for _, notes := range []int{100, 1000, 5000} {
		b.Run(fmt.Sprintf("notes=%d", notes), func(b *testing.B) {
			benchVault(b, notes)

			s, err := NewServer(b.Context())
			require.NoError(b, err)

			src := "See [["
			require.NoError(b, os.WriteFile("source.md", []byte(src), 0o644))

			params := &protocol.CompletionParams{
				TextDocumentPositionParams: protocol.TextDocumentPositionParams{
					TextDocument: protocol.TextDocumentIdentifier{URI: "file://" + filepath.Join(mustGetwd(b), "source.md")},
					Position:     protocol.Position{Line: 0, Character: protocol.UInteger(len(src))},
				},
			}

			b.ResetTimer()

			for b.Loop() {
				_, err := s.TextDocumentCompletion(nil, params)
				require.NoError(b, err)
			}
		})
	}
}

func mustGetwd(b *testing.B) string {
	b.Helper()

	wd, err := os.Getwd()
	require.NoError(b, err)

	return wd
}
