package note

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"os/exec"
	"path/filepath"
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
//
// TODO (jamesl33): Defer reading the note body, until required?
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

// Name - TODO
func (n *Note) Name() string {
	return strings.TrimSuffix(filepath.Base(n.Path), ".md")
}

// Checksum - TODO
func (n *Note) Checksum() (uint32, error) {
	hasher := crc32.NewIEEE()

	err := n.WriteTo(hasher)
	if err != nil {
		return 0, fmt.Errorf("%w", err) // TODO
	}

	return hasher.Sum32(), nil
}

// Edit - TODO
func (n *Note) Edit(ctx context.Context) error {
	// TODO
	ed := os.Getenv("EDITOR")

	// TODO
	if ed == "" {
		return errors.New("no editor set") // TODO
	}

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

	// Re-read the note
	r, err := New(n.Path)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// Shallow copy
	*n = *r

	return nil
}

// Write - TODO
func (n *Note) Write() error {
	file, err := os.OpenFile(n.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}
	defer file.Close()

	err = n.WriteTo(file)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// WriteTo - TODO
func (n *Note) WriteTo(w io.Writer) error {
	// marker - TODO
	const marker = "---\n"

	b := bufio.NewWriter(w)

	_, err := b.WriteString(marker)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = yaml.NewEncoder(b).Encode(n.Frontmatter)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	_, err = b.WriteString(marker)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	_, err = b.WriteString(n.Body)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = b.Flush()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// Links - TODO
func (n *Note) Links() []string {
	var (
		// TODO
		re = regexp.MustCompile(`\[\[(?P<link>.*?)(\|(?P<text>.*?))?\]\]`)

		// TODO
		matches = re.FindAllStringSubmatch(n.Body, -1)
	)

	links := make([]string, 0)

	for _, match := range matches {
		links = append(links, match[re.SubexpIndex("link")])
	}

	return links
}

// String0 - TODO
func (n *Note) String0() string {
	str := fmt.Sprintf(
		"%s\x01%s\x01%s\x01%s",
		icolor.Blue(filepath.Dir(n.Path)),
		icolor.Yellow(n.Frontmatter.Title),
		icolor.Cyan(strings.Join(n.Frontmatter.Tags, ",")),
		n.Path,
	)

	return str
}
