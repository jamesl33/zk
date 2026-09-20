package matcher

import (
	"fmt"

	"github.com/jamesl33/zk/internal/note"
	"go.yaml.in/yaml/v4"
)

// Frontmatter returns a matcher which matches against the marshalled representation of the frontmatter.
func Frontmatter(f, g, r string) (Matcher, error) {
	extract := func(n *note.Note) (string, error) {
		data, err := yaml.Marshal(n.Frontmatter)
		if err != nil {
			return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
		}

		return string(data), nil
	}

	return text(f, g, r, func(n *note.Note) (string, error) { return extract(n) })
}
