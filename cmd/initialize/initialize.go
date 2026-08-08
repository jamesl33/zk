package initialize

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed .geminiignore
var ignore []byte

//go:embed .gemini
var settings embed.FS

//go:embed GEMINI.md
var instructions []byte

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
		Short: "Sets up Gemini CLI with instructions/settings on how to interact with the Zettelkasten",
		Use:   "initialize",
		RunE:  func(cmd *cobra.Command, _ []string) error { return index.Run(cmd.Context()) },
	}

	return &cmd
}

// Run initialization.
func (i *Initialize) Run(ctx context.Context) error {
	err := os.RemoveAll(".gemini")
	if err != nil {
		return fmt.Errorf("failed to remove existing '.gemini' directory: %w", err)
	}

	err = os.WriteFile(".geminiignore", ignore, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write '.geminiignore': %w", err)
	}

	err = os.WriteFile("GEMINI.md", instructions, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write 'GEMINI.md': %w", err)
	}

	err = fs.WalkDir(settings, ".", i.cp)
	if err != nil {
		return fmt.Errorf("failed to walk '.gemini' directory: %w", err)
	}

	return nil
}

// cp the given file into the '.gemini' directory on disk.
func (i *Initialize) cp(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if entry.IsDir() {
		return nil
	}

	err = os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	data, err := settings.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read embedded file %s: %w", path, err)
	}

	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
