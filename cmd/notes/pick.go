package notes

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// PickOptions defines the options for the pick command.
type PickOptions struct{}

// Pick defines the struct for the pick command.
type Pick struct {
	PickOptions
}

// NewPick creates a new command for picking a note using 'fzf'.
func NewPick() *cobra.Command {
	var pick Pick

	cmd := cobra.Command{
		Short: "Pick a note using 'fzf', supports the output from 'zk'",
		Use:   "pick",
		RunE:  func(cmd *cobra.Command, _ []string) error { return pick.Run(cmd.Context()) },
	}

	return &cmd
}

// Run the note picker.
func (p *Pick) Run(ctx context.Context) error {
	var buffer bytes.Buffer

	cmd := exec.CommandContext(
		ctx,
		"fzf",
		"--ansi",
		"--exit-0",
		"--select-1",
		`--delimiter=\x01`,
		"--with-nth={1} {2} [{3}]",
		`--preview=bat --color=always --style=numbers {4}`,
		"--tac",
	)

	// We must pass all these through
	cmd.Stdin = os.Stdin
	cmd.Stdout = &buffer
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		return p.item(buffer)
	}

	// TODO (jamesl33): There's no entries; not a fan of the exit status though.
	// TODO (jamesl33): User has exited.
	if cmd.ProcessState != nil && (cmd.ProcessState.ExitCode() == 1 || cmd.ProcessState.ExitCode() == 130) {
		return nil
	}

	return fmt.Errorf("failed to pick note: %w", err)
}

// item prints the path to the chosen note.
func (p *Pick) item(buffer bytes.Buffer) error {
	var (
		// Reverse the printing for note lines
		split = bytes.Split(buffer.Bytes(), []byte{0x01})
		// The result is empty
		empty = len(split) == 0
		// The result is a single newline
		nl = len(split) == 1 && bytes.Equal(split[0], []byte("\n"))
	)

	// If empty, don't print anything
	if empty || nl {
		return nil
	}

	fmt.Printf("%s", split[len(split)-1])

	return nil
}
