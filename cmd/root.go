package cmd

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/jamesl33/zk/cmd/index"
	"github.com/jamesl33/zk/cmd/note"
	"github.com/jamesl33/zk/cmd/notes"
	"github.com/jamesl33/zk/cmd/tags"
	"github.com/spf13/cobra"
)

// rootCommand - TODO
var rootCommand = &cobra.Command{
	// TODO
	Short: "",
	// TODO
	Long: "",
	// TODO
	SilenceErrors: true,
	// TODO
	SilenceUsage: true,
	// TODO
	TraverseChildren: true,
}

// init - TODO
func init() {
	rootCommand.AddCommand(
		// TODO
		index.NewIndex(),
		// TODO
		note.NewNote(),
		// TODO
		notes.NewNotes(),
		// TODO
		tags.NewTags(),
	)
}

// Execute - TODO
func Execute() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := rootCommand.ExecuteContext(ctx)
	if err == nil {
		return nil
	}

	// TODO
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return fmt.Errorf("%w", err) // TODO
}
