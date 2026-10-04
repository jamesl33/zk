// Package initialize provides the 'zk initialize' command, which wires Claude Code up to work
// with a Zettelkasten.
package initialize

import (
	"github.com/jamesl33/zk/cmd/initialize/claude"
	"github.com/spf13/cobra"
)

// NewInitialize creates a new command for initializing Claude Code for use with 'zk'.
func NewInitialize() *cobra.Command {
	cmd := cobra.Command{
		Short: "Sets up Claude Code to work with the Zettelkasten",
		Use:   "initialize",
	}

	cmd.AddCommand(
		claude.NewClaude(),
	)

	return &cmd
}
