package note

import (
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

	// Body is the - front-matter excluded - note body.
	Body string `json:"body,omitempty" jsonschema:"The notes content/body, excluding frontmatter"`
}

// New returns a new note.
func New(path string) (*Note, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read note at %q: %w", path, err)
	}

	const marker = "---\n"

	var (
		sfm = bytes.Index(data, []byte(marker))
		efm = bytes.Index(data[sfm+1:], []byte(marker))
	)

	var fm Frontmatter

	err = yaml.Unmarshal(data[sfm:efm], &fm)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	note := Note{
		Path:        path,
		Frontmatter: fm,
		Body:        string(data[sfm+efm+len(marker)+1:]),
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
		matches = regex.Link.FindAllStringSubmatch(n.Body, -1)
		links   = make([]string, 0)
	)

	for _, match := range matches {
		links = append(links, match[regex.Link.SubexpIndex("link")])
	}

	return links
}

// String returns a string representation of the note.
func (n *Note) String() string {
	var builder strings.Builder

	// TODO (jamesl33): Fine to just ignore this?
	_, _ = n.WriteTo(&builder)

	return builder.String()
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
