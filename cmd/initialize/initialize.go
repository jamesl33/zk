package initialize

import (
	"context"
	"fmt"
	"os"

	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed settings.json
var settings []byte

//go:embed instructions.md
var instructions []byte

//go:embed researcher.md
var researcher []byte

//go:embed maintainer.md
var maintainer []byte

// InitializeOptions defines the options for the initialize command.
type InitializeOptions struct{}

// Initialize defines the struct for the initialize command.
type Initialize struct {
	InitializeOptions
}

// NewInitialize creates a new command for initializing Gemini for use with 'zk'.
func NewInitialize() *cobra.Command {
	var index Initialize

	cmd := cobra.Command{
		Short: "Creates a 'GEMINI.md' file, with instructions on how to interact with the Zettelkasten",
		Use:   "initialize",
		RunE:  func(cmd *cobra.Command, _ []string) error { return index.Run(cmd.Context()) },
	}

	return &cmd
}

// Run initialization.
func (i *Initialize) Run(ctx context.Context) error {
	err := os.MkdirAll(".gemini", 0o755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = os.WriteFile(".gemini/settings.json", settings, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	err = os.WriteFile("GEMINI.md", instructions, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write instructions: %w", err)
	}

	err = os.MkdirAll(".gemini/skills/researcher", 0o755)
	if err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	err = os.WriteFile(".gemini/skills/researcher/SKILL.md", researcher, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write researcher skill: %w", err)
	}

	err = os.MkdirAll(".gemini/skills/maintainer", 0o755)
	if err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	err = os.WriteFile(".gemini/skills/maintainer/SKILL.md", maintainer, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write maintainer skill: %w", err)
	}

	return nil
}
