package tags

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/note"
	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v4"
)

// GenerateOptions - TODO
type GenerateOptions struct{}

// Generate - TODO
type Generate struct {
	GenerateOptions
}

// NewGenerate - TODO
func NewGenerate() *cobra.Command {
	var generate Generate

	cmd := cobra.Command{
		// TODO
		Short: "",
		// TODO
		Use: "generate",
		// TODO
		Args: cobra.MaximumNArgs(1),
		// TODO
		RunE: func(cmd *cobra.Command, args []string) error { return generate.Run(cmd.Context(), args) },
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
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return fmt.Errorf("failed to create ollama client: %w", err)
	}

	example := "```yaml\ntags:\n  - tag_1\n  - tag_2\n```"

	prompt := `

%s

Using the above context, output up to five "tags" to help catagorize this note. Use the following format.

%s

You must use lower-case and only output tags using the snake case style.

Don't use tags unless there's enough information to catagorize.`

	// TODO (jamesl33): Note body should remove wiki-links (replace with the title?).
	req := api.GenerateRequest{
		Model:  "gemma3:4b",
		Prompt: fmt.Sprintf(prompt, n.Body, example),
	}

	var o strings.Builder

	err = client.Generate(ctx, &req, func(resp api.GenerateResponse) error { o.WriteString(resp.Response); return nil })
	if err != nil {
		return fmt.Errorf("failed to generate tags: %w", err)
	}

	// Extracts the YAML from the markdown code-block
	re := regexp.MustCompile(`\x60\x60\x60yaml(?P<tags>[\S\s]*?.*)\x60\x60\x60`)

	// Extract the tags
	m := re.FindStringSubmatch(o.String())

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

	n.Frontmatter.Tags = overlay.Tags

	err = n.Write()
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	return nil
}
