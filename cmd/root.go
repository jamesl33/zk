package cmd

import (
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
	return rootCommand.Execute()
}
