package matcher

import (
	"slices"

	"github.com/jamesl33/zk/internal/note"
)

// Tagged returns a matcher which looks for a given tag.
func Tagged(tag string) Matcher {
	return func(n *note.Note) bool { return slices.Contains(n.Frontmatter.Tags, tag) }
}
