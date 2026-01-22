package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/jamesl33/zk/cmd/mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

// MCPOptions defines the options for the mcp command.
type MCPOptions struct{}

// MCP defines the struct for the mcp command.
type MCP struct {
	MCPOptions
}

// NewMCP creates a new command for acting as a model context protocol server.
func NewMCP() *cobra.Command {
	var mcp MCP

	cmd := cobra.Command{
		Short:  "Runs as a model context protocol server",
		Use:    "mcp",
		RunE:   func(cmd *cobra.Command, _ []string) error { return mcp.Run(cmd.Context()) },
		Hidden: true,
	}

	return &cmd
}

// Run configures and executes the model context protocol server. It initializes the server implementation, registers
// available tools with their descriptions, and then starts the server using stdio for transport.
func (m *MCP) Run(ctx context.Context) error {
	impl := mcp.Implementation{
		Name:    "zk",
		Version: "v0.1.0",
	}

	server := mcp.NewServer(&impl, nil)

	description := `
List notes in a directory.

This tool returns all the notes within a directory so is therefore useful when you want to obtain *all* the notes relating to a topic.

There's no filtering, so it's best for obtaining all the notes you wish to read.
`

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "list_notes", Description: strings.TrimSpace(description)},
		tools.ListNotes,
	)

	description = `
Search the title, tags and content of notes using a regular expression.

It's a case sensitive, multi-line search by default; to make it case insensive, use "(?i:$EXPRESSION)".

This is an excellent way to find notes, it should be used when:

	- Searching for a specific title
	- Searching for a tag, which is a lower case, snake case identifier (e.g. linux, thru_hiking)
	- Searching for words or phrases which may be within a note

The full note content isn't returned, you must read it using 'read_file'.
`

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "regex_search_notes", Description: strings.TrimSpace(description)},
		tools.RegexSearchNotes,
	)

	description = `
Search for notes similar to the given phrase, topic or idea.

This is an excellent way to find notes, it should be used when:

	- Searching for notes which are similar, or related to a topic (but perhaps don't mention it directly)

The full note content isn't returned, you must read it using 'read_file'.
`

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "semantic_search_notes", Description: strings.TrimSpace(description)},
		tools.SemanticSearchNotes,
	)

	description = `
Find notes that contain information that's semantically similar or link to/from a given note.

This is a great way to increase the amount of context you have before performing a task, it should be used when:

	- You have a note and you'd like to find other notes that are similar (vector search)
	- You have a note and you'd like to find other notes that are linked to/from the note (direct mentions)
	- You have a 'bibliographic' note, and wish to find quotes/citations from the book

The full note content isn't returned, you must read it using 'read_file'.
	`

	mcp.AddTool(
		server,
		&mcp.Tool{Name: "find_related_notes", Description: strings.TrimSpace(description)},
		tools.FindRelatedNotes,
	)

	err := server.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}
