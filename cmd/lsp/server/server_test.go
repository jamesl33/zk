package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// chdir changes into the given directory for the duration of the test.
func chdir(t *testing.T, dir string) {
	t.Helper()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })
}

func definitionParams(t *testing.T, path string, line, char int) *protocol.DefinitionParams {
	t.Helper()

	abs, err := filepath.Abs(path)
	require.NoError(t, err)

	return &protocol.DefinitionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: "file://" + abs},
			Position:     protocol.Position{Line: protocol.UInteger(line), Character: protocol.UInteger(char)},
		},
	}
}

func referenceParams(t *testing.T, path string) *protocol.ReferenceParams {
	t.Helper()

	abs, err := filepath.Abs(path)
	require.NoError(t, err)

	return &protocol.ReferenceParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: "file://" + abs},
		},
	}
}

func TestTextDocumentReferencesNoBacklinks(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	target := note.Note{Path: "20060102150405.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Target"}}
	require.NoError(t, target.Write())

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentReferences(nil, referenceParams(t, "20060102150405.md"))
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestTextDocumentReferencesSingleBacklink(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	target := note.Note{Path: "20060102150405.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Target"}}
	require.NoError(t, target.Write())

	source := note.Note{Path: "20060102150406.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Source"}}
	source.SetBody("See [[20060102150405|Target]]")
	require.NoError(t, source.Write())

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentReferences(nil, referenceParams(t, "20060102150405.md"))
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Contains(t, result[0].URI, "20060102150406.md")
}

func TestTextDocumentReferencesMultipleBacklinks(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	target := note.Note{Path: "20060102150405.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Target"}}
	require.NoError(t, target.Write())

	first := note.Note{Path: "20060102150406.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "First"}}
	first.SetBody("See [[20060102150405|Target]]")
	require.NoError(t, first.Write())

	second := note.Note{Path: "20060102150407.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Second"}}
	second.SetBody("Also see [[20060102150405|Target]]")
	require.NoError(t, second.Write())

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentReferences(nil, referenceParams(t, "20060102150405.md"))
	require.NoError(t, err)
	require.Len(t, result, 2)

	uris := []string{result[0].URI, result[1].URI}
	assert.Contains(t, uris, "file://"+mustAbs(t, "20060102150406.md"))
	assert.Contains(t, uris, "file://"+mustAbs(t, "20060102150407.md"))
}

func TestTextDocumentReferencesNotANote(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))
	require.NoError(t, os.WriteFile("source.md", []byte("No frontmatter here"), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentReferences(nil, referenceParams(t, "source.md"))
	require.NoError(t, err)
	assert.Nil(t, result)
}

func mustAbs(t *testing.T, path string) string {
	t.Helper()

	abs, err := filepath.Abs(path)
	require.NoError(t, err)

	return abs
}

func TestTextDocumentDefinitionSelectsLinkUnderCursor(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	target := note.Note{Path: "20060102150406.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Target"}}
	require.NoError(t, target.Write())

	other := note.Note{Path: "20060102150405.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Other"}}
	require.NoError(t, other.Write())

	src := "See [[20060102150405|Other]] and [[20060102150406|Target]]"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	// Position the cursor within the second link, which points at "target".
	char := strings.Index(src, "20060102150406")

	result, err := s.TextDocumentDefinition(nil, definitionParams(t, "source.md", 0, char))
	require.NoError(t, err)

	loc, ok := result.(protocol.Location)
	require.True(t, ok)
	assert.Contains(t, loc.URI, "20060102150406.md")
}

func TestTextDocumentDefinitionNoLinkAtCursor(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	src := "No links on this line"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentDefinition(nil, definitionParams(t, "source.md", 0, 0))
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestTextDocumentDefinitionBrokenLink(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	src := "See [[20060102150406|Missing]]"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	char := strings.Index(src, "20060102150406")

	result, err := s.TextDocumentDefinition(nil, definitionParams(t, "source.md", 0, char))
	require.NoError(t, err)
	assert.Nil(t, result)
}
