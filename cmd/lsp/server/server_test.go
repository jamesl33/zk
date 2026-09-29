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
	require.NoError(t, target.Create())

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
	require.NoError(t, target.Create())

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
	require.NoError(t, target.Create())

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
	require.NoError(t, target.Create())

	other := note.Note{Path: "20060102150405.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Other"}}
	require.NoError(t, other.Create())

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
	assert.Equal(t, loc.Range.Start, loc.Range.End)
}

func TestTextDocumentDefinitionEscapesURI(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))
	require.NoError(t, os.Mkdir("1 Projects", 0o755))

	target := note.Note{Path: "1 Projects/20060102150406.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Target"}}
	require.NoError(t, target.Create())

	src := "See [[20060102150406]]"
	require.NoError(t, os.WriteFile("source.md", []byte(src), 0o644))

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	result, err := s.TextDocumentDefinition(nil, definitionParams(t, "source.md", 0, strings.Index(src, "2006")))
	require.NoError(t, err)

	loc, ok := result.(protocol.Location)
	require.True(t, ok)
	assert.Contains(t, loc.URI, "1%20Projects/20060102150406.md")
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

func TestTextDocumentDidChangeCompletesUnsavedLink(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))

	target := note.Note{Path: "20060102150405.md", Frontmatter: note.Frontmatter{Type: "permanent", Title: "Target"}}
	require.NoError(t, target.Create())

	// The file on disk has no open link, the editor buffer does.
	require.NoError(t, os.WriteFile("source.md", []byte("See"), 0o644))

	abs, err := filepath.Abs("source.md")
	require.NoError(t, err)

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	ctx := &glsp.Context{Notify: func(string, any) {}}

	err = s.TextDocumentDidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{URI: "file://" + abs},
		},
		ContentChanges: []any{protocol.TextDocumentContentChangeEventWhole{Text: "See [["}},
	})
	require.NoError(t, err)

	result, err := s.TextDocumentCompletion(nil, completionParams(t, "source.md", 0, len("See [[")))
	require.NoError(t, err)
	assert.Len(t, result, 1)

	// Once closed, the disk is read again.
	err = s.TextDocumentDidClose(nil, &protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{URI: "file://" + abs},
	})
	require.NoError(t, err)

	result, err = s.TextDocumentCompletion(nil, completionParams(t, "source.md", 0, len("See [[")))
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestTextDocumentDidChangePublishesDiagnostics(t *testing.T) {
	tmp := t.TempDir()
	chdir(t, tmp)

	require.NoError(t, os.Mkdir(".zk", 0o755))
	require.NoError(t, os.WriteFile("source.md", []byte("nothing"), 0o644))

	abs, err := filepath.Abs("source.md")
	require.NoError(t, err)

	s, err := NewServer(t.Context())
	require.NoError(t, err)

	var published protocol.PublishDiagnosticsParams

	ctx := &glsp.Context{
		Notify: func(_ string, params any) { published = params.(protocol.PublishDiagnosticsParams) },
	}

	err = s.TextDocumentDidChange(ctx, &protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: protocol.TextDocumentIdentifier{URI: "file://" + abs},
		},
		ContentChanges: []any{protocol.TextDocumentContentChangeEventWhole{Text: "See [[20060102150405]]"}},
	})
	require.NoError(t, err)

	assert.Len(t, published.Diagnostics, 1)
}

func TestUTF16Offsets(t *testing.T) {
	// '😀' is 4 bytes but 2 UTF-16 units, 'é' is 2 bytes but 1 unit.
	line := "é😀[[x]]"

	assert.Equal(t, 2, byteOffset(line, 1))
	assert.Equal(t, 6, byteOffset(line, 3))
	assert.Equal(t, len(line), byteOffset(line, 100))

	assert.Equal(t, 1, utf16Offset(line, 2))
	assert.Equal(t, 3, utf16Offset(line, 6))
}

func TestFindLinksReportsUTF16Columns(t *testing.T) {
	links := findLinks("😀 [[20060102150405]]")
	require.Len(t, links, 1)

	assert.Equal(t, 3, links[0].StartChar)
	assert.Equal(t, 3+len("[[20060102150405]]"), links[0].EndChar)
}

func TestLinkAtCursorUTF16(t *testing.T) {
	lines := []string{"😀 [[20060102150405]]"}

	// Column 3 is the first bracket in UTF-16 units, but byte offset 5.
	assert.Equal(t, "20060102150405", linkAtCursor(lines, protocol.Position{Line: 0, Character: 3}))
	assert.Empty(t, linkAtCursor(lines, protocol.Position{Line: 0, Character: 1}))
}
