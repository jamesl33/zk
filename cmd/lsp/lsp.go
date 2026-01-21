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

// LSPOptions - TODO
type LSPOptions struct{}

// LSP - TODO
type LSP struct {
	LSPOptions
}

// NewLSP - TODO
func NewLSP() *cobra.Command {
	var lsp LSP

	cmd := cobra.Command{
		Short: "Starts an LSP server for zk",
		Use:   "lsp",
		RunE:  func(cmd *cobra.Command, _ []string) error { return lsp.Run(cmd.Context()) },
	}

	return &cmd
}

// Run - TODO
func (l *LSP) Run(ctx context.Context) error {
	svr, err := lspserver.NewServer(ctx)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return server.NewServer(svr, "zk", false).RunStdio()
}
