package note

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/jamesl33/zk/internal/ptr"
	"github.com/jamesl33/zk/internal/regex"
	"github.com/mattn/go-isatty"
	"go.yaml.in/yaml/v4"
)

func init() {
	color.NoColor = false
}

// Note is a markdown note.
type Note struct {
	// Path to the note.
	Path string `json:"path" jsonschema:"The relative path to the note on disk"`

	// Frontmatter metadata for the note.
	Frontmatter Frontmatter `json:"frontmatter" jsonschema:"The notes YAML frontmatter"`

	// body is the lazy loaded body of the note.
	body *string
}

// New returns a new note.
//
// TODO (jamesl33): Improve the handling of non-note markdown files.
func New(path string) (*Note, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open note at %q: %w", path, err)
	}
	defer file.Close()

	// Create a scanner, which we'll use to pluck the frontmatter
	scanner := bufio.NewScanner(file)

	// Add the frontmatter scanner
	scanner.Split(scan)

	ok := scanner.Scan()
	if !ok {
		return nil, fmt.Errorf("failed to discard the first frontmatter marker: %w", scanner.Err())
	}

	ok = scanner.Scan()
	if !ok {
		return nil, fmt.Errorf("failed to scan frontmatter: %w", scanner.Err())
	}

	var fm Frontmatter

	err = yaml.Unmarshal(scanner.Bytes(), &fm)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	note := Note{
		Path:        path,
		Frontmatter: fm,
	}

	return &note, nil
}

// Name returns the notes name.
func (n *Note) Name() string {
	return strings.TrimSuffix(filepath.Base(n.Path), ".md")
}

// GetBody returns the note body.
func (n *Note) GetBody() (string, error) {
	if n.body != nil {
		return *n.body, nil
	}

	file, err := os.Open(n.Path)
	if err != nil {
		return "", fmt.Errorf("failed to open note at %q: %w", n.Path, err)
	}
	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get note stats: %w", err)
	}

	// Create a scanner, which we'll use to pluck the body
	scanner := bufio.NewScanner(file)

	// Add the frontmatter scanner
	scanner.Split(scan)

	// Set the buffer to the size of the file, so we can read the whole thing in one go
	scanner.Buffer(make([]byte, stats.Size()), int(stats.Size()))

	// Throw away the first marker
	ok := scanner.Scan()
	if !ok {
		return "", fmt.Errorf("failed to discard the first frontmatter marker: %w", scanner.Err())
	}

	// Throw away the frontmatter
	ok = scanner.Scan()
	if !ok {
		return "", fmt.Errorf("failed to discard the frontmatter: %w", scanner.Err())
	}

	// Scan the body
	ok = scanner.Scan()
	if !ok {
		return "", fmt.Errorf("failed to scan the frontmatter: %w", scanner.Err())
	}

	n.body = ptr.To(scanner.Text())

	return *n.body, nil
}

// SetBody overwrites the note body.
//
// NOTE: This is in-memory only, use 'Write' to persist the changes.
func (n *Note) SetBody(body string) {
	n.body = &body
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

	body, err := n.GetBody()
	if err != nil {
		return 0, fmt.Errorf("failed to get body: %w", err)
	}

	_, err = b.WriteString(body)
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
func (n *Note) Links() ([]string, error) {
	body, err := n.GetBody()
	if err != nil {
		return nil, fmt.Errorf("failed to get body: %w", err)
	}

	var (
		matches = regex.Link.FindAllStringSubmatch(body, -1)
		links   = make([]string, 0)
	)

	for _, match := range matches {
		links = append(links, match[regex.Link.SubexpIndex("link")])
	}

	return links, nil
}

// Text returns the full note text.
func (n *Note) Text() (string, error) {
	var builder strings.Builder

	_, err := n.WriteTo(&builder)
	if err != nil {
		return "", fmt.Errorf("failed to write note: %w", err)
	}

	return builder.String(), nil
}

// String0 returns a null-delimited representation of the note, useful for "picking" (i.e. 'fzf').
//
// TODO (jamesl33): Think of a better name for this.
func (n *Note) String0() string {
	var (
		blue   = color.New(color.FgBlue).SprintFunc()
		dir    = blue(filepath.Dir(n.Path))
		yellow = color.New(color.FgYellow).SprintFunc()
		title  = yellow(n.Frontmatter.Title)
		cyan   = color.New(color.FgCyan).SprintFunc()
		tags   = cyan(strings.Join(n.Frontmatter.Tags, ","))
	)

	if isatty.IsTerminal(os.Stdout.Fd()) {
		return fmt.Sprintf("%s %s [%s] %s", dir, title, tags, n.Path)
	}

	return fmt.Sprintf("%s\x01%s\x01[%s]\x01%s", dir, title, tags, n.Path)
}
