package note

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
		Short: "",
		// TODO
		Use: "update",
		// TODO
		Args: cobra.ExactArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return update.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run - TODO
//
// TODO (jamesl33): We should update indexes post-update.
func (u *Update) Run(ctx context.Context, args []string) error {
	path, err := u.path(args)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	//
	// TODO (jamesl33): Return an error if no editor is setup?
	ed := os.Getenv("EDITOR")

	// TODO
	cmd := exec.CommandContext(
		ctx,
		ed,
		strings.TrimSuffix(path, "\n"),
	)

	// TODO
	cmd.Stdin = os.Stdin

	// TODO
	cmd.Stdout = os.Stdout

	// TODO
	cmd.Stderr = os.Stderr

	// TODO
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// path - TODO
func (u *Update) path(args []string) (string, error) {
	// TODO
	if args[0] != "-" {
		return args[0], nil
	}

	reader := bufio.NewReader(os.Stdin)

	path, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("%w", err) // TODO
	}

	// TODO
	path = strings.TrimSuffix(path, "\n")

	return path, nil
}
