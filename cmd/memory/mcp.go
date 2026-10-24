package memory

import (
	"context"

	"github.com/jamesl33/zk/cmd/mcp"
	"github.com/spf13/cobra"
)

// MCPOptions defines the options for the mcp command.
type MCPOptions struct{}

// MCP defines the struct for the mcp command.
type MCP struct {
	MCPOptions
}

// NewMCP creates a new command for serving the memory vault as a model context protocol server.
func NewMCP() *cobra.Command {
	var server MCP

	cmd := cobra.Command{
		Short: "Runs the model context protocol server against the memory vault",
		Use:   "mcp",
		RunE:  func(cmd *cobra.Command, _ []string) error { return server.Run(cmd.Context()) },
	}

	return &cmd
}

// Run the server. Entering the vault here (rather than passing '--vault') creates it first, as the server may start
// before the SessionStart hook has run.
func (m *MCP) Run(ctx context.Context) error {
	_, err := enter()
	if err != nil {
		return err
	}

	var server mcp.MCP

	return server.Run(ctx)
}
