package notes

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// PickOptions - TODO
type PickOptions struct{}

// Pick - TODO
type Pick struct {
	PickOptions
}

// NewPick - TODO
func NewPick() *cobra.Command {
	var pick Pick

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "pick",
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return pick.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run - TODO
func (l *Pick) Run(ctx context.Context, args []string) error {
	var buffer bytes.Buffer

	// TODO
	cmd := exec.CommandContext(
		ctx,
		"fzf",
		"--ansi",
		"--select-1",
		"--with-nth={1} ({2})",
		`--delimiter=\x00`,
	)

	// TODO
	cmd.Stdin = os.Stdin

	// TODO
	cmd.Stdout = &buffer

	// TODO
	cmd.Stderr = os.Stderr

	// TODO
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	split := bytes.Split(buffer.Bytes(), []byte{0})

	// TODO
	if len(split) == 0 {
		return nil
	}

	// TODO
	fmt.Printf("%s", split[len(split)-1])

	return nil
}
