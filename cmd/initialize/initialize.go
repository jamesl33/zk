// Package initialize provides the 'zk initialize' command, which wires an AI coding agent up to work
// with a Zettelkasten via 'zk mcp'.
package initialize

import (
	"github.com/jamesl33/zk/cmd/initialize/claude"
	"github.com/jamesl33/zk/cmd/initialize/gemini"
	"github.com/spf13/cobra"
)

// NewInitialize creates a new command for initializing an AI coding agent for use with 'zk'.
func NewInitialize() *cobra.Command {
	cmd := cobra.Command{
		Short: "Sets up an AI coding agent with instructions/settings on how to interact with the Zettelkasten",
		Use:   "initialize",
	}

	cmd.AddCommand(
		claude.NewClaude(),
		gemini.NewGemini(),
	)

	return &cmd
}
