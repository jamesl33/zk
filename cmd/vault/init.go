// Package vault provides the 'zk init' command, which creates a new Zettelkasten vault.
package vault

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// directories is the PARA layout created in every new vault.
var directories = []string{
	"0 Inbox",
	"1 Projects",
	"2 Areas",
	"3 Resources",
	"4 Archives",
	"5 Bibliography",
}

// InitOptions defines the options for the init command.
type InitOptions struct{}

// Init defines the struct for the init command.
type Init struct {
	InitOptions
}

// NewInit creates a new command for creating a vault.
func NewInit() *cobra.Command {
	var v Init

	cmd := cobra.Command{
		Short: "Creates a vault (the '.zk' directory and PARA layout), leaving anything that exists untouched",
		Use:   "init [directory]",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return v.Run(".")
			}

			return v.Run(args[0])
		},
	}

	return &cmd
}

// Run creates the vault in the given directory, which is created if required.
func (i *Init) Run(path string) error {
	for _, dir := range append([]string{".zk"}, directories...) {
		err := os.MkdirAll(filepath.Join(path, dir), 0o755)
		if err != nil {
			return fmt.Errorf("failed to create '%s': %w", dir, err)
		}
	}

	return nil
}
