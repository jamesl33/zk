package note

import (
	"bytes"
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

	"github.com/fatih/color"
	"go.yaml.in/yaml/v4"
)

func init() {
	color.NoColor = false
}

// Note is a markdown note.
type Note struct {
	// Path to the note.
	Path string

	// Frontmatter metadata for the note.
	Frontmatter Frontmatter

	// Body is the - front-matter excluded - note body.
	Body string
}

// New returns a new note.
//
// TODO (jamesl33): Defer reading the note body, until required?
func New(path string) (*Note, error) {
	re := regexp.MustCompile(`^---[\S\s]*?---\n.*`)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read note at %q: %w", path, err)
	}

	loc := re.FindIndex(data)

	// We require the offset/length
	if len(loc) != 2 {
		return nil, errors.New("failed to extract front-matter") // TODO (jamesl33): Better error for this?
	}

	var (
		offset = int64(loc[0])
		length = int64(loc[1])
	)

	var fm Frontmatter

	err = yaml.Unmarshal(data[offset:length], &fm)
	if err != nil {
		return nil, fmt.Errorf("failed to parse front-matter: %w", err)
	}

	note := Note{
		Path:        path,
		Frontmatter: fm,
		Body:        string(data[length:]),
	}

	return &note, nil
}

// Name returns the notes name.
func (n *Note) Name() string {
	return strings.TrimSuffix(filepath.Base(n.Path), ".md")
}

// Checksum returns a checksum of the entire note (including front-matter).
func (n *Note) Checksum() (uint32, error) {
	hasher := crc32.NewIEEE()

	_, err := n.WriteTo(hasher)
	if err != nil {
		return 0, fmt.Errorf("failed to hash note: %w", err)
	}

	return hasher.Sum32(), nil
}

// Edit opens the note in the users default editor.
func (n *Note) Edit(ctx context.Context) error {
	ed := os.Getenv("EDITOR")

	if ed == "" {
		return fmt.Errorf("no editor set in the %q environment variable", "EDITOR")
	}

	cmd := exec.CommandContext(
		ctx,
		ed,
		strings.TrimSuffix(n.Path, "\n"),
	)

	// We must pass all these through
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	// Re-read the note
	r, err := New(n.Path)
	if err != nil {
		return fmt.Errorf("failed to read updated note: %w", err)
	}

	// Shallow copy
	*n = *r

	return nil
}

// Write the note out to disk.
func (n *Note) Write() error {
	file, err := os.OpenFile(n.Path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open file at %q: %w", n.Path, err)
	}
	defer file.Close()

	_, err = n.WriteTo(file)
	if err != nil {
		return fmt.Errorf("failed to write note to file: %w", err)
	}

	return nil
}

// WriteTo writes the note out to the given writer.
func (n *Note) WriteTo(w io.Writer) (int64, error) {
	// marker used for the YAML front-matter.
	const marker = "---\n"

	var b bytes.Buffer

	_, err := b.WriteString(marker)
	if err != nil {
		return 0, fmt.Errorf("failed to write first marker: %w", err)
	}

	err = yaml.NewEncoder(&b).Encode(n.Frontmatter)
	if err != nil {
		return 0, fmt.Errorf("failed to write front-matter: %w", err)
	}

	_, err = b.WriteString(marker)
	if err != nil {
		return 0, fmt.Errorf("failed to write second marker: %w", err)
	}

	_, err = b.WriteString(n.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to write body: %w", err)
	}

	bw, err := io.Copy(w, bytes.NewReader(b.Bytes()))
	if err != nil {
		return bw, fmt.Errorf("failed to write to: %w", err)
	}

	return bw, nil
}

// Links returns the names of other notes mentioned in this note.
func (n *Note) Links() []string {
	var (
		// Captures wiki-style links
		re = regexp.MustCompile(`\[\[(?P<link>.*?)(\|(?P<text>.*?))?\]\]`)

		// All the matches within the note body
		matches = re.FindAllStringSubmatch(n.Body, -1)
	)

	links := make([]string, 0)

	for _, match := range matches {
		links = append(links, match[re.SubexpIndex("link")])
	}

	return links
}

// String0 returns a null-delimited representation of the note, useful for "picking" (i.e. 'fzf').
func (n *Note) String0() string {
	var (
		yellow = color.New(color.FgYellow).SprintFunc()
		blue   = color.New(color.FgBlue).SprintFunc()
		cyan   = color.New(color.FgCyan).SprintFunc()
	)

	str := fmt.Sprintf(
		"%s\x01%s\x01%s\x01%s",
		blue(filepath.Dir(n.Path)),
		yellow(n.Frontmatter.Title),
		cyan(strings.Join(n.Frontmatter.Tags, ",")),
		n.Path,
	)

	return str
}
