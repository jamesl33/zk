package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/tliron/glsp"
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

func completionParams(t *testing.T, path string, line, char int) *protocol.CompletionParams {
	t.Helper()

	abs, err := filepath.Abs(path)
	require.NoError(t, err)

	return &protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: "file://" + abs},
			Position:     protocol.Position{Line: protocol.UInteger(line), Character: protocol.UInteger(char)},
		},
	}
}

func TestTextDocumentCompletionInsideOpenLink(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	target := note.Note{Path: "20060102150405.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Target"}}
	require.NoError(t, target.Write())

	src := "See [["
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentCompletion(nil, completionParams(t, "source.md", 0, len(src)))
	require.NoError(t, err)

	items, ok := result.([]protocol.CompletionItem)
	require.True(t, ok)
	require.Len(t, items, 1)
	assert.Contains(t, items[0].Label, "Target")
}

func TestTextDocumentCompletionOutsideLink(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	src := "No link here"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentCompletion(nil, completionParams(t, "source.md", 0, len(src)))
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestTextDocumentCompletionAfterClosedLink(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	src := "See [[20060102150405|Target]] "
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentCompletion(nil, completionParams(t, "source.md", 0, len(src)))
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestTextDocumentDidSavePublishesBrokenLinkDiagnostic(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	src := "See [[20060102150405|Missing]]"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	abs, err := filepath.Abs("source.md")
	require.NoError(t, err)

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	var published protocol.PublishDiagnosticsParams

	ctx := &glsp.Context{
		Notify: func(method string, params any) {
			assert.Equal(t, string(protocol.ServerTextDocumentPublishDiagnostics), method)
			published = params.(protocol.PublishDiagnosticsParams)
		},
	}

	err = s.TextDocumentDidSave(ctx, &protocol.DidSaveTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: "file://" + abs},
	})
	require.NoError(t, err)

	require.Len(t, published.Diagnostics, 1)
	assert.Contains(t, published.Diagnostics[0].Message, "20060102150405")
}

func TestTextDocumentDidOpenNoDiagnosticsForValidLink(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	target := note.Note{Path: "20060102150405.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Target"}}
	require.NoError(t, target.Write())

	src := "See [[20060102150405|Target]]"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	abs, err := filepath.Abs("source.md")
	require.NoError(t, err)

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	var published protocol.PublishDiagnosticsParams

	ctx := &glsp.Context{
		Notify: func(method string, params any) {
			published = params.(protocol.PublishDiagnosticsParams)
		},
	}

	err = s.TextDocumentDidOpen(ctx, &protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{URI: "file://" + abs},
	})
	require.NoError(t, err)

	assert.Empty(t, published.Diagnostics)
}

func hoverParams(t *testing.T, path string, line, char int) *protocol.HoverParams {
	t.Helper()

	abs, err := filepath.Abs(path)
	require.NoError(t, err)

	return &protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{URI: "file://" + abs},
			Position:     protocol.Position{Line: protocol.UInteger(line), Character: protocol.UInteger(char)},
		},
	}
}

func TestTextDocumentHoverShowsLinkedNoteTitle(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	target := note.Note{Path: "20060102150405.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Target"}}
	require.NoError(t, target.Write())

	src := "See [[20060102150405|Target]]"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	char := strings.Index(src, "20060102150405")

	result, err := s.TextDocumentHover(nil, hoverParams(t, "source.md", 0, char))
	require.NoError(t, err)
	require.NotNil(t, result)

	content, ok := result.Contents.(protocol.MarkupContent)
	require.True(t, ok)
	assert.Contains(t, content.Value, "Target")
}

func TestTextDocumentHoverNoLinkAtCursor(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	src := "No links on this line"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentHover(nil, hoverParams(t, "source.md", 0, 0))
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestTextDocumentHoverBrokenLink(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	src := "See [[20060102150406|Missing]]"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	char := strings.Index(src, "20060102150406")

	result, err := s.TextDocumentHover(nil, hoverParams(t, "source.md", 0, char))
	require.NoError(t, err)
	assert.Nil(t, result)
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
