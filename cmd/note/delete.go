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
	"github.com/jamesl33/zk/internal/vault"
	"github.com/spf13/cobra"
)

// DeleteOptions defines the options for the delete command.
type DeleteOptions struct{}

// Delete defines the struct for the delete command.
type Delete struct {
	DeleteOptions
}

// NewDelete creates a new command for deleting a note.
func NewDelete() *cobra.Command {
	var del Delete

	cmd := cobra.Command{
		Short: "Delete a note",
		Use:   "delete [path]",
		RunE:  func(cmd *cobra.Command, args []string) error { return del.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run deletes the given note.
func (d *Delete) Run(_ context.Context, args []string) error {
	if _, err := vault.RootAbs("."); err != nil {
		return err
	}

	path, err := d.path(args)

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

	err = os.Remove(n.Path)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	return nil
}

// path returns the path to the target note.
func (d *Delete) path(args []string) (string, error) {
	// Path provided, use that
	if len(args) != 0 && args[0] != "-" {
		return args[0], nil
	}

	// No path provided, read from stdin
	path, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("failed to read from stdin: %w", err)
	}

	// Strip whitespace
	path = strings.TrimSuffix(path, "\n")

	// No data at all, treat this the same as user-cancellation
	if path == "" {
		return "", io.EOF
	}

	return path, nil
}
