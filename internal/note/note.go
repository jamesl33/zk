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
	re := regexp.MustCompile(`^---[\S\s]*?---.*`)

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

// String0 - TODO
func (n Note) String0() string {
	return fmt.Sprintf("%s\x00%s", icolor.Yellow(n.Frontmatter.Title), icolor.Blue(n.Path))
}
