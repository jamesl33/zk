// Package claude sets up Claude Code to work with 'zk'.
package claude

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/jamesl33/zk/cmd/initialize/assets"
	"github.com/spf13/cobra"
)

const (
	instructionsFile = "ZK.md"

	pluginInstructions = `Install the 'zk' plugin (skills and MCP server) from within Claude Code:

  /plugin marketplace add jamesl33/zk
  /plugin install zk@zk`
)

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
		Short: "Writes ZK.md instructions, imports them into CLAUDE.md and explains how to install the plugin",
		Use:   "claude",
		RunE:  func(cmd *cobra.Command, _ []string) error { return claude.Run(cmd.OutOrStdout()) },
	}

	return &cmd
}

// Run initialization.
func (c *Claude) Run(out io.Writer) error {
	err := os.WriteFile(instructionsFile, assets.Instructions, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write '%s': %w", instructionsFile, err)
	}

	err = ensureImport("CLAUDE.md", "@"+instructionsFile)
	if err != nil {
		return fmt.Errorf("failed to import '%s' into 'CLAUDE.md': %w", instructionsFile, err)
	}

	fmt.Fprintln(out, pluginInstructions)

	return nil
}

// ensureImport prepends 'line' to the file at 'path', creating it if required, unless it is already present.
func ensureImport(path, line string) error {
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	for existing := range strings.Lines(string(data)) {
		if strings.TrimSpace(existing) == line {
			return nil
		}
	}

	content := line + "\n"
	if len(data) > 0 {
		content += "\n" + string(data)
	}

	return os.WriteFile(path, []byte(content), 0o644)
}
