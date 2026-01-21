package server

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/ptr"
	"github.com/jamesl33/zk/internal/regex"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// Server defines the LSP server.
//
// TODO (jamesl33): Add support for finding references (e.g. all the notes that link to a note).
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
	}

	return &server, nil
}

// Initialize the LSP server's capabilities.
func (s *Server) Initialize(_ *glsp.Context, _ *protocol.InitializeParams) (any, error) {
	capabilities := s.CreateServerCapabilities()

	capabilities.DefinitionProvider = true

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
		line     = lines[params.Position.Line]
		matches  = regex.Link.FindStringSubmatch(line)
		expected = 1 + regex.Link.NumSubexp()
	)

	if len(matches) != expected {
		return nil, fmt.Errorf("unexpected number of regex matches")
	}

	name := matches[regex.Link.SubexpIndex("link")]

	l, err := lister.NewLister(
		// TODO (jamesl33): This should probably be 'git rev-parse --show-toplevel'?
		lister.WithPath("."),
		lister.WithMatcher(matcher.Name(name)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	dst, err := l.One(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	// Note not found, broken link?
	if dst == nil {
		return nil, nil
	}

	abs, err := filepath.Abs(dst.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path of destination note: %w", err)
	}

	// TODO (jamesl33): Jump directly to the title?
	rng := protocol.Range{
		Start: protocol.Position{Line: 0, Character: 0},
		End:   protocol.Position{Line: 0, Character: 0},
	}

	loc := protocol.Location{
		URI:   "file://" + abs,
		Range: rng,
	}

	return loc, nil
}
