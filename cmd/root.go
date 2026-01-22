package cmd

import (
	"context"
	"errors"
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
	"github.com/spf13/cobra"
)

// rootCommand defines the root of the command chain.
var rootCommand = &cobra.Command{
	Short:            "A composable command-line tool for interacting with a Markdown Zettelkasten.",
	Use:              "zk",
	SilenceErrors:    true,
	SilenceUsage:     true,
	TraverseChildren: true,
}

// init sets up the CLI.
func init() {
	rootCommand.AddCommand(
		initialize.NewInitialize(),
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
