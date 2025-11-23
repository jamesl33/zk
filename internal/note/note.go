package note

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v4"
)

// Note - TODO
type Note struct {
	// Path - TODO
	Path string
}

// NewNote - TODO
func NewNote(path string) *Note {
	note := Note{
		Path: path,
	}

	return &note
}

// Name - TODO
func (n *Note) Name() string {
	return strings.TrimSuffix(filepath.Base(n.Path), ".md")
}

// Frontmatter - TODO
//
// TODO (jamesl33): This should be optimised.
func (n *Note) Frontmatter() (Frontmatter, error) {
	re := regexp.MustCompile(`^---[\S\s]*?---`)

	file, err := os.Open(n.Path)
	if err != nil {
		return Frontmatter{}, fmt.Errorf("%w", err) // TODO
	}
	defer file.Close()

	loc := re.FindReaderIndex(bufio.NewReader(file))

	// TODO
	if len(loc) != 2 {
		return Frontmatter{}, errors.New("not found") // TODO
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return Frontmatter{}, fmt.Errorf("%w", err) // TODO
	}

	var (
		offset = int64(loc[0])
		length = int64(loc[1])
		s      = io.NewSectionReader(file, offset, length)
	)

	var p Frontmatter

	err = yaml.NewDecoder(s).Decode(&p)
	if err != nil {
		return Frontmatter{}, fmt.Errorf("%w", err) // TODO
	}

	return p, nil
}

// Body - TODO
//
// TODO (jamesl33): This should be optimised.
func (n *Note) Body() (string, error) {
	re := regexp.MustCompile(`^---[\S\s]*?---`)

	file, err := os.Open(n.Path)
	if err != nil {
		return "", fmt.Errorf("%w", err) // TODO
	}
	defer file.Close()

	loc := re.FindReaderIndex(bufio.NewReader(file))

	// TODO
	if len(loc) != 2 {
		return "", errors.New("not found") // TODO
	}

	offset := int64(loc[1])

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("%w", err) // TODO
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%w", err) // TODO
	}

	return string(data), nil
}
