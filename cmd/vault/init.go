// Package vault provides the 'zk init' command, which creates a new Zettelkasten vault.
package vault

import (
	"github.com/jamesl33/zk/internal/vault"
	"github.com/spf13/cobra"
)

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
	return vault.Init(path)
}
