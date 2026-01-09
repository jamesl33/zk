package cmd

import (
	"context"
	"errors"
	"os/signal"
	"syscall"

	"github.com/jamesl33/zk/cmd/initialize"
	"github.com/jamesl33/zk/cmd/note"
	"github.com/jamesl33/zk/cmd/notes"
	"github.com/jamesl33/zk/cmd/tags"
	"github.com/spf13/cobra"
)

// rootCommand - TODO
var rootCommand = &cobra.Command{
	// TODO
	Short: "A composable command-line tool for interacting with a Markdown Zettelkasten.",
	// TODO
	Long: "",
	// TODO
	SilenceErrors: true,
	// TODO
	SilenceUsage: true,
	// TODO
	TraverseChildren: true,
}

// init sets up the CLI.
func init() {
	rootCommand.AddCommand(initialize.NewInitialize(), note.NewNote(), notes.NewNotes(), tags.NewTags())
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
