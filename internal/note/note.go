package note

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	icolor "github.com/jamesl33/zk/internal/color"
	"go.yaml.in/yaml/v4"
)

// Note - TODO
type Note struct {
	// Path - TODO
	Path string

	// Frontmatter - TODO
	Frontmatter Frontmatter

	// Body - TODO
	Body string
}

// New - TODO
func New(path string) (*Note, error) {
	re := regexp.MustCompile(`^---[\S\s]*?---\n.*`)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	loc := re.FindIndex(data)

	// TODO
	if len(loc) != 2 {
		return nil, errors.New("not found") // TODO
	}

	var (
		offset = int64(loc[0])
		length = int64(loc[1])
	)

	var fm Frontmatter

	err = yaml.Unmarshal(data[offset:length], &fm)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	note := Note{
		Path:        path,
		Frontmatter: fm,
		Body:        string(data[length:]),
	}

	return &note, nil
}

// Edit - TODO
//
// TODO (jamesl33): This should re-read the new after exit.
func (n Note) Edit(ctx context.Context) error {
	// TODO
	//
	// TODO (jamesl33): Return an error if no editor is setup?
	ed := os.Getenv("EDITOR")

	// TODO
	cmd := exec.CommandContext(
		ctx,
		ed,
		strings.TrimSuffix(n.Path, "\n"),
	)

	// TODO
	cmd.Stdin = os.Stdin

	// TODO
	cmd.Stdout = os.Stdout

	// TODO
	cmd.Stderr = os.Stderr

	// TODO
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// Write - TODO
func (n Note) Write() error {
	file, err := os.OpenFile(n.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}
	defer file.Close()

	_, err = file.WriteString("---\n")
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = yaml.NewEncoder(file).Encode(n.Frontmatter)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	_, err = file.WriteString("---\n")
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	_, err = file.WriteString(n.Body)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// String0 - TODO
func (n Note) String0() string {
	return fmt.Sprintf("%s\x00%s", icolor.Yellow(n.Frontmatter.Title), icolor.Blue(n.Path))
}
