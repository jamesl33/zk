package note

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jamesl33/zk/internal/note"
	"github.com/spf13/cobra"
)

// UpdateOptions - TODO
type UpdateOptions struct{}

// Update - TODO
type Update struct {
	UpdateOptions
}

// NewUpdate - TODO
func NewUpdate() *cobra.Command {
	var update Update

	cmd := cobra.Command{
		// TODO
		Short: "Open a note in the default text editor",
		// TODO
		Use: "update",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return update.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run opens a new editor for the given note.
func (u *Update) Run(ctx context.Context, args []string) error {
	path, err := u.path(args)

	// User didn't provide input, exit cleanly (this better handles exiting pickers early)
	if errors.Is(err, io.EOF) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to get path: %w", err)
	}

	n, err := note.New(path)
	if err != nil {
		return fmt.Errorf("failed to open note: %w", err)
	}

	err = n.Edit(ctx)
	if err != nil {
		return fmt.Errorf("failed to edit note: %w", err)
	}

	return nil
}

// path returns the path to the target note.
func (u *Update) path(args []string) (string, error) {
	// Path provided, use that
	if len(args) != 0 && args[0] != "-" {
		return args[0], nil
	}

	// No path provided, read from stdin
	path, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("%w", err) // TODO
	}

	// Strip whitespace
	path = strings.TrimSuffix(path, "\n")

	return path, nil
}
