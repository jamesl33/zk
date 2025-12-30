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
		RunE: func(cmd *cobra.Command, _ []string) error { return pick.Run(cmd.Context()) },
	}

	return &cmd
}

// Run the note picker.
//
// TODO (jamesl33): Better handle the case where no elements are piped into 'fzf'.
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
		`--preview=zk note summarize {4}`,
		"--preview-window=wrap",
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
	split := bytes.Split(buffer.Bytes(), []byte{0x01})

	if len(split) == 0 {
		return nil
	}

	fmt.Printf("%s", split[len(split)-1])

	return nil
}
