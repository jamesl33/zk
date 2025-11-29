package tags

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/lister"
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

// Run - TODO
func (g *Generate) Run(ctx context.Context, args []string) error {
	path := "."

	if len(args) >= 1 {
		path = args[0]
	}

	lister, err := lister.NewLister(
		lister.WithPath(path),
	)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	err = iterator.ForEach2(lister.Many(ctx), func(n *note.Note) error {
		return g.generate(ctx, n)
	})
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}

// generate - TODO
func (g *Generate) generate(ctx context.Context, n *note.Note) error {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
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
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	re := regexp.MustCompile(`\x60\x60\x60yaml(?P<tags>[\S\s]*?.*)\x60\x60\x60`)

	// TODO
	m := re.FindStringSubmatch(o.String())

	// TODO
	if len(m) != 2 {
		return nil
	}

	// TODO
	tags := m[re.SubexpIndex("tags")]

	// TODO
	var overlay struct {
		Tags []string `yaml:"tags"`
	}

	err = yaml.Unmarshal([]byte(tags), &overlay)
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	// TODO
	for i := range overlay.Tags {
		// TODO
		overlay.Tags[i] = strings.ReplaceAll(overlay.Tags[i], " ", "_")

		// TODO
		overlay.Tags[i] = strings.ReplaceAll(overlay.Tags[i], "-", "_")
	}

	n.Frontmatter.Tags = overlay.Tags

	err = n.Write()
	if err != nil {
		return fmt.Errorf("%w", err) // TODO
	}

	return nil
}
