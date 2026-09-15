// Package claude sets up Claude Code to work with 'zk'.
package claude

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/jamesl33/zk/cmd/initialize/assets"
	"github.com/spf13/cobra"
)

//go:embed mcp.json
var mcpConfig []byte

//go:embed settings.json
var settings []byte

// ClaudeOptions defines the options for the initialize claude command.
type ClaudeOptions struct{}

// Claude defines the struct for the initialize claude command.
type Claude struct {
	ClaudeOptions
}

// NewClaude creates a new command for initializing Claude Code for use with 'zk'.
func NewClaude() *cobra.Command {
	var claude Claude

	cmd := cobra.Command{
		Short: "Sets up Claude Code with instructions/settings on how to interact with the Zettelkasten",
		Use:   "claude",
		RunE:  func(cmd *cobra.Command, _ []string) error { return claude.Run(cmd.Context()) },
	}

	return &cmd
}

// Run initialization.
func (c *Claude) Run(_ context.Context) error {
	err := os.RemoveAll(".claude")
	if err != nil {
		return fmt.Errorf("failed to remove existing '.claude' directory: %w", err)
	}

	err = os.WriteFile("CLAUDE.md", assets.Instructions, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write 'CLAUDE.md': %w", err)
	}

	err = os.WriteFile(".mcp.json", mcpConfig, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write '.mcp.json': %w", err)
	}

	err = os.MkdirAll(".claude", 0o755)
	if err != nil {
		return fmt.Errorf("failed to create '.claude' directory: %w", err)
	}

	err = os.WriteFile(".claude/settings.json", settings, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write '.claude/settings.json': %w", err)
	}

	err = assets.CopySkills(".claude/skills")
	if err != nil {
		return fmt.Errorf("failed to write '.claude/skills': %w", err)
	}

	return nil
}
