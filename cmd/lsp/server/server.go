package server

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/notes"
	"github.com/jamesl33/zk/internal/ptr"
	"github.com/jamesl33/zk/internal/regex"
	"github.com/jamesl33/zk/internal/vault"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Server defines the LSP server.
type Server struct {
	protocol.Handler
	ctx context.Context
}

// NewServer creates a new LSP server.
func NewServer(ctx context.Context) (*Server, error) {
	server := Server{
		ctx: ctx,
	}

	server.Handler = protocol.Handler{
		Initialize:             server.Initialize,
		Initialized:            server.Initialized,
		Shutdown:               server.Shutdown,
		SetTrace:               server.SetTrace,
		TextDocumentDefinition: server.TextDocumentDefinition,
		TextDocumentReferences: server.TextDocumentReferences,
		TextDocumentCompletion: server.TextDocumentCompletion,
		TextDocumentDidOpen:    server.TextDocumentDidOpen,
		TextDocumentDidSave:    server.TextDocumentDidSave,
	}

	return &server, nil
}

// Initialize the LSP server's capabilities.
func (s *Server) Initialize(_ *glsp.Context, _ *protocol.InitializeParams) (any, error) {
	capabilities := s.CreateServerCapabilities()

	capabilities.DefinitionProvider = true
	capabilities.ReferencesProvider = true
	capabilities.CompletionProvider = &protocol.CompletionOptions{TriggerCharacters: []string{"["}}

	si := protocol.InitializeResultServerInfo{
		Name:    "zk",
		Version: ptr.To("0.1.0"),
	}

	result := protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo:   &si,
	}

	return result, nil
}

// Initialized is called after the client has been initialized.
func (s *Server) Initialized(_ *glsp.Context, _ *protocol.InitializedParams) error {
	return nil
}

// Shutdown is called when the client requests to shut down the server.
func (s *Server) Shutdown(_ *glsp.Context) error {
	return nil
}

// SetTrace sets the trace value.
func (s *Server) SetTrace(_ *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)

	return nil
}

// TextDocumentDefinition provides the definition for a symbol at a given position.
func (s *Server) TextDocumentDefinition(_ *glsp.Context, params *protocol.DefinitionParams) (any, error) {
	u, err := url.Parse(params.TextDocument.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document URI: %w", err)
	}

	src, err := os.ReadFile(u.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %w", err)
	}

	lines := strings.Split(string(src), "\n")

	if params.Position.Line >= uint32(len(lines)) {
		return nil, nil
	}

	var (
		cur     = lines[params.Position.Line]
		matches = regex.Link.FindAllStringSubmatchIndex(cur, -1)
		linkIdx = regex.Link.SubexpIndex("link")
	)

	name := ""

	for _, match := range matches {
		start, end := match[0], match[1]

		if int(params.Position.Character) < start || int(params.Position.Character) >= end {
			continue
		}

		name = cur[match[2*linkIdx]:match[2*linkIdx+1]]

		break
	}

	// The cursor isn't positioned on a link.
	if name == "" {
		return nil, nil
	}

	root, err := vault.Root(".")
	if err != nil {
		return nil, fmt.Errorf("failed to find vault root: %w", err)
	}

	l, err := lister.NewLister(
		lister.WithPath(root),
		lister.WithMatcher(matcher.Name(name)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	dst, err := l.One(s.ctx)

	// Note not found, broken link?
	if errors.Is(err, lister.ErrNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	abs, err := filepath.Abs(dst.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of destination note: %w", err)
	}

	body, err := dst.Text()
	if err != nil {
		return nil, fmt.Errorf("failed to get note text: %w", err)
	}

	lines = strings.Split(body, "\n")

	// Place the cursor at the beginning of the note title
	var (
		idx  = slices.IndexFunc(lines, func(s string) bool { return strings.HasPrefix(s, "title: ") })
		line = max(0, idx)
		char = len("title: ")
	)

	rng := protocol.Range{
		Start: protocol.Position{Line: protocol.UInteger(line), Character: protocol.UInteger(char)},
		End:   protocol.Position{Line: 0, Character: 0},
	}

	loc := protocol.Location{
		URI:   "file://" + abs,
		Range: rng,
	}

	return loc, nil
}

// TextDocumentReferences provides the locations of all notes which link to the note in the given document.
func (s *Server) TextDocumentReferences(_ *glsp.Context, params *protocol.ReferenceParams) ([]protocol.Location, error) {
	u, err := url.Parse(params.TextDocument.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document URI: %w", err)
	}

	n, err := note.New(u.Path)

	// The document isn't a note, nothing links to it.
	if errors.Is(err, note.ErrNotNote) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open note: %w", err)
	}

	locs := []protocol.Location{}

	err = notes.LinkedTo(s.ctx, n, func(src *note.Note) {
		abs, aerr := filepath.Abs(src.Path)
		if aerr != nil {
			return
		}

		body, berr := src.Text()
		if berr != nil {
			return
		}

		for _, m := range findLinks(body) {
			if m.Name != n.Name() {
				continue
			}

			locs = append(locs, protocol.Location{
				URI: "file://" + abs,
				Range: protocol.Range{
					Start: protocol.Position{Line: protocol.UInteger(m.Line), Character: protocol.UInteger(m.Start)},
					End:   protocol.Position{Line: protocol.UInteger(m.Line), Character: protocol.UInteger(m.End)},
				},
			})
		}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find linked notes: %w", err)
	}

	return locs, nil
}

// TextDocumentCompletion provides note name completions inside an open WikiLink.
func (s *Server) TextDocumentCompletion(_ *glsp.Context, params *protocol.CompletionParams) (any, error) {
	u, err := url.Parse(params.TextDocument.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document URI: %w", err)
	}

	src, err := os.ReadFile(u.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %w", err)
	}

	lines := strings.Split(string(src), "\n")

	if params.Position.Line >= uint32(len(lines)) {
		return nil, nil
	}

	var (
		cur  = lines[params.Position.Line]
		char = min(int(params.Position.Character), len(cur))
		open = strings.LastIndex(cur[:char], "[[")
	)

	// The cursor isn't inside an open link.
	if open == -1 || strings.LastIndex(cur[:char], "]]") > open {
		return nil, nil
	}

	root, err := vault.Root(".")
	if err != nil {
		return nil, fmt.Errorf("failed to find vault root: %w", err)
	}

	l, err := lister.NewLister(
		lister.WithPath(root),
		lister.WithMatcher(matcher.Any()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	items := []protocol.CompletionItem{}

	err = iterator.ForEach2(l.Many(s.ctx), hs.Infallible(func(n *note.Note) {
		items = append(items, protocol.CompletionItem{
			Label:      fmt.Sprintf("%s %s", n.Name(), n.Frontmatter.Title),
			InsertText: ptr.To(fmt.Sprintf("%s|%s", n.Name(), n.Frontmatter.Title)),
			FilterText: ptr.To(fmt.Sprintf("%s %s", n.Name(), n.Frontmatter.Title)),
			Detail:     ptr.To(n.Path),
		})
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	return items, nil
}

// TextDocumentDidOpen publishes diagnostics for the note that was opened.
func (s *Server) TextDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.publishDiagnostics(ctx, params.TextDocument.URI)

	return nil
}

// TextDocumentDidSave publishes diagnostics for the note that was saved.
func (s *Server) TextDocumentDidSave(ctx *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	s.publishDiagnostics(ctx, params.TextDocument.URI)

	return nil
}

// publishDiagnostics scans the note at the given URI for broken links and notifies the client.
func (s *Server) publishDiagnostics(ctx *glsp.Context, uri protocol.DocumentUri) {
	diags, err := s.diagnostics(uri)
	if err != nil {
		return
	}

	ctx.Notify(string(protocol.ServerTextDocumentPublishDiagnostics), protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diags,
	})
}

// diagnostics scans the note at the given URI for links which point at notes that don't exist.
func (s *Server) diagnostics(uri protocol.DocumentUri) ([]protocol.Diagnostic, error) {
	u, err := url.Parse(string(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to parse document URI: %w", err)
	}

	src, err := os.ReadFile(u.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read source file: %w", err)
	}

	root, err := vault.Root(".")
	if err != nil {
		return nil, fmt.Errorf("failed to find vault root: %w", err)
	}

	diags := []protocol.Diagnostic{}

	for _, m := range findLinks(string(src)) {
		exists, err := noteExists(s.ctx, root, m.Name)
		if err != nil {
			return nil, err
		}

		if exists {
			continue
		}

		diags = append(diags, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: protocol.UInteger(m.Line), Character: protocol.UInteger(m.Start)},
				End:   protocol.Position{Line: protocol.UInteger(m.Line), Character: protocol.UInteger(m.End)},
			},
			Severity: ptr.To(protocol.DiagnosticSeverityWarning),
			Source:   ptr.To("zk"),
			Message:  fmt.Sprintf("note not found: %q", m.Name),
		})
	}

	return diags, nil
}

// noteExists reports whether a note with the given name exists in the vault rooted at root.
func noteExists(ctx context.Context, root, name string) (bool, error) {
	l, err := lister.NewLister(
		lister.WithPath(root),
		lister.WithMatcher(matcher.Name(name)),
	)
	if err != nil {
		return false, fmt.Errorf("failed to create lister: %w", err)
	}

	_, err = l.One(ctx)

	if errors.Is(err, lister.ErrNotFound) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to get note: %w", err)
	}

	return true, nil
}

// linkMatch is a single WikiLink occurrence within a note body.
type linkMatch struct {
	Line       int
	Start, End int
	Name       string
}

// findLinks returns every WikiLink occurrence within the given body.
func findLinks(body string) []linkMatch {
	var (
		matches = []linkMatch{}
		linkIdx = regex.Link.SubexpIndex("link")
	)

	for i, line := range strings.Split(body, "\n") {
		for _, m := range regex.Link.FindAllStringSubmatchIndex(line, -1) {
			matches = append(matches, linkMatch{
				Line:  i,
				Start: m[0],
				End:   m[1],
				Name:  line[m[2*linkIdx]:m[2*linkIdx+1]],
			})
		}
	}

	return matches
}
