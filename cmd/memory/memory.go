// Package memory provides the hidden 'zk memory' commands, which an agent's hooks call to back its persistent memory
// with a zk vault.
package memory

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/jamesl33/zk/internal/memory"
	"github.com/spf13/cobra"
)

// NewMemory creates a new command for backing an agent's persistent memory with a vault.
func NewMemory() *cobra.Command {
	cmd := cobra.Command{
		Short:  "Backs an agent's persistent memory with a vault",
		Use:    "memory",
		Hidden: true,
	}

	cmd.AddCommand(
		NewHook(),
		NewMCP(),
		NewSave(),
	)

	return &cmd
}

// enter ensures the memory vault and runs in it, returning its path.
func enter() (string, error) {
	vault, err := memory.VaultPath()
	if err != nil {
		return "", err
	}

	err = memory.EnsureVault(vault)
	if err != nil {
		return "", fmt.Errorf("failed to ensure vault: %w", err)
	}

	err = os.Chdir(vault)
	if err != nil {
		return "", fmt.Errorf("failed to enter vault: %w", err)
	}

	return vault, nil
}

// write the hook output.
func write(w io.Writer, out memory.Output) error {
	return json.NewEncoder(w).Encode(out)
}
