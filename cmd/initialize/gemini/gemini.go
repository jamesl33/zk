// Package gemini sets up the Gemini CLI to work with 'zk'.
package gemini

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/jamesl33/zk/cmd/initialize/assets"
	"github.com/spf13/cobra"
)

//go:embed settings.json
var settings []byte

//go:embed policies/zk.toml
var policy []byte

// GeminiOptions defines the options for the initialize gemini command.
type GeminiOptions struct{}

// Gemini defines the struct for the initialize gemini command.
type Gemini struct {
	GeminiOptions
}

// NewGemini creates a new command for initializing Gemini CLI for use with 'zk'.
func NewGemini() *cobra.Command {
	var gemini Gemini

	cmd := cobra.Command{
		Short: "Sets up Gemini CLI with instructions/settings on how to interact with the Zettelkasten",
		Use:   "gemini",
		RunE:  func(cmd *cobra.Command, _ []string) error { return gemini.Run(cmd.Context()) },
	}

	return &cmd
}

// Run initialization.
func (g *Gemini) Run(_ context.Context) error {
	err := os.RemoveAll(".gemini")
	if err != nil {
		return fmt.Errorf("failed to remove existing '.gemini' directory: %w", err)
	}

	err = os.WriteFile("GEMINI.md", assets.Instructions, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write 'GEMINI.md': %w", err)
	}

	err = os.MkdirAll(".gemini/policies", 0o755)
	if err != nil {
		return fmt.Errorf("failed to create '.gemini/policies' directory: %w", err)
	}

	err = os.WriteFile(".gemini/settings.json", settings, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write '.gemini/settings.json': %w", err)
	}

	err = os.WriteFile(".gemini/policies/zk.toml", policy, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write '.gemini/policies/zk.toml': %w", err)
	}

	err = assets.CopySkills(".gemini/skills")
	if err != nil {
		return fmt.Errorf("failed to write '.gemini/skills': %w", err)
	}

	return nil
}
