package tags

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jamesl33/zk/internal/ai"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/links"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/ptr"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

// GenerateOptions defines the options for the generate command.
type GenerateOptions struct{}

// Generate defines the struct for the generate command.
type Generate struct {
	GenerateOptions
}

// NewGenerate creates a new command for generating tags.
func NewGenerate() *cobra.Command {
	var generate Generate

	cmd := cobra.Command{
		Short: "Generate tags for notes, based on the note content",
		Use:   "generate [directory | path]",
		Args:  cobra.MaximumNArgs(1),
		RunE:  func(cmd *cobra.Command, args []string) error { return generate.Run(cmd.Context(), args) },
	}

	return &cmd
}

// Run tag generation.
func (g *Generate) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lister.Many(ctx), func(n *note.Note) error {
		return g.generate(ctx, n)
	})
	if err != nil {
		return fmt.Errorf("failed to update notes: %w", err)
	}

	return nil
}

// generate tags for the given note.
func (g *Generate) generate(ctx context.Context, n *note.Note) error {
	// Read the body before creating a copy of the note
	body, err := n.GetBody()
	if err != nil {
		return fmt.Errorf("failed to get body: %w", err)
	}

	if len(body) == 0 {
		return nil
	}

	// Create a copy of the note, that we can modify in-place
	cp := ptr.To(*n)

	err = links.Replace(ctx, cp)
	if err != nil {
		return fmt.Errorf("failed to replace links: %w", err)
	}

	client, err := ai.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	example := "```yaml\ntags:\n  - tag_1\n  - tag_2\n```"

	prompt := `

%s

Using the above context, output up to five "tags" to help catagorize this note. Use the following format.

%s

You must use lower-case and only output tags using the snake case style.

Don't use tags unless there's enough information to catagorize.`

	prompt = fmt.Sprintf(prompt, body, example)

	content, err := client.Generate(ctx, prompt)
	if err != nil {
		return fmt.Errorf("failed to generate tags: %w", err)
	}

	// Extracts the YAML from the markdown code-block
	re := regexp.MustCompile(`\x60\x60\x60yaml(?P<tags>[\S\s]*?.*)\x60\x60\x60`)

	// Extract the tags
	m := re.FindStringSubmatch(content)

	// We didn't find everything, ignore
	if len(m) != 2 {
		return nil
	}

	// Extract the tags
	tags := m[re.SubexpIndex("tags")]

	// overlay allows extracting the tags
	var overlay struct {
		Tags []string `yaml:"tags"`
	}

	err = yaml.Unmarshal([]byte(tags), &overlay)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tags: %w", err)
	}

	for i := range overlay.Tags {
		// Coerce spaces into snake case
		overlay.Tags[i] = strings.ReplaceAll(overlay.Tags[i], " ", "_")

		// Coerce kebab casing into snake case
		overlay.Tags[i] = strings.ReplaceAll(overlay.Tags[i], "-", "_")
	}

	// Update the original note so the unchanged body is re-written to disk
	n.Frontmatter.Tags = overlay.Tags

	err = n.Write()
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	return nil
}
