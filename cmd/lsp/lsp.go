package lsp

import (
	"context"
	"fmt"

	lspserver "github.com/jamesl33/zk/cmd/lsp/server"
	"github.com/spf13/cobra"
	"github.com/tliron/glsp/server"

	// Must include a backend implementation.
	_ "github.com/tliron/commonlog/simple"
)

// LSPOptions defines the options for the lsp command.
type LSPOptions struct{}

// LSP defines the struct for the lsp command.
type LSP struct {
	LSPOptions
}

// NewLSP creates a new command for starting the zk LSP server.
func NewLSP() *cobra.Command {
	var lsp LSP

	cmd := cobra.Command{
		Short:  "Starts an LSP server for zk",
		Use:    "lsp",
		RunE:   func(cmd *cobra.Command, _ []string) error { return lsp.Run(cmd.Context()) },
		Hidden: true,
	}

	return &cmd
}

// Run runs the LSP server.
func (l *LSP) Run(ctx context.Context) error {
	svr, err := lspserver.NewServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to create LSP server: %w", err)
	}

	return server.NewServer(svr, "zk", false).RunStdio()
}
