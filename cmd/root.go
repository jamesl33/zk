package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jamesl33/zk/cmd/initialize"
	"github.com/jamesl33/zk/cmd/links"
	"github.com/jamesl33/zk/cmd/lint"
	"github.com/jamesl33/zk/cmd/lsp"
	"github.com/jamesl33/zk/cmd/mcp"
	"github.com/jamesl33/zk/cmd/note"
	"github.com/jamesl33/zk/cmd/notes"
	"github.com/jamesl33/zk/cmd/tags"
	"github.com/jamesl33/zk/cmd/vault"
	"github.com/spf13/cobra"
)

// root is the directory to run in, in place of the current working directory.
var root string

// rootCommand defines the root of the command chain.
var rootCommand = &cobra.Command{
	Short:            "A composable command-line tool for interacting with a Markdown Zettelkasten.",
	Use:              "zk",
	SilenceErrors:    true,
	SilenceUsage:     true,
	TraverseChildren: true,
	PersistentPreRunE: func(*cobra.Command, []string) error {
		return enterVault(root)
	},
}

// enterVault changes into the given directory so that every command resolves the vault from there, a no-op if empty.
func enterVault(path string) error {
	if path == "" {
		return nil
	}

	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("failed to enter vault: %w", err)
	}

	return nil
}

// init sets up the CLI.
func init() {
	rootCommand.PersistentFlags().StringVar(
		&root,
		"vault",
		os.Getenv("ZK_VAULT"),
		"Run in this directory rather than the current one (env: ZK_VAULT)",
	)

	rootCommand.AddCommand(
		initialize.NewInitialize(),
		vault.NewInit(),
		links.NewLinks(),
		lint.NewLint(),
		lsp.NewLSP(),
		mcp.NewMCP(),
		note.NewNote(),
		notes.NewNotes(),
		tags.NewTags(),
	)
}

// Execute 'zk'.
func Execute() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := rootCommand.ExecuteContext(ctx)
	if err == nil {
		return nil
	}

	// The user canceled, don't output an error (useful for piping)
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err // Purposefully not wrapped
}
