package note

import (
	"errors"
	"fmt"
	"os"
	"regexp"

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

// NewNote - TODO
func NewNote(path string) (*Note, error) {
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
