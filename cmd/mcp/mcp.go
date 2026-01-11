package mcp

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/cmd/mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

// MCPOptions - TODO
type MCPOptions struct {
	// Title - TODO
	Title string
}

// MCP - TODO
type MCP struct {
	MCPOptions
}

// NewMCP - TODO
func NewMCP() *cobra.Command {
	var mcp MCP

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "mcp",
		// TODO
		RunE: func(cmd *cobra.Command, _ []string) error { return mcp.Run(cmd.Context()) },
	}

	return &cmd
}

// Run creates a new bibliographic note.
func (m *MCP) Run(ctx context.Context) error {
	impl := mcp.Implementation{
		Name:    "zk",
		Version: "v0.1.0",
	}

	server := mcp.NewServer(&impl, nil)

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "list_notes", Description: "List notes in a directory"},
		tools.ListNotes,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "regex_search_notes", Description: "Search note titles, tags and body using a regular expression"},
		tools.RegexSearchNotes,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "semantic_search_notes", Description: "Search for notes which are semantically similar to a query, phrase or word"},
		tools.SemanticSearchNotes,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "find_related_notes", Description: "Performs a semantic search, to find related notes"},
		tools.FindRelatedNotes,
	)

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "read_note", Description: "Read a note from disk"},
		tools.ReadNote,
	)

	err := server.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}
