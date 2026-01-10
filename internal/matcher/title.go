package matcher

import "github.com/jamesl33/zk/internal/note"

// Title returns a matcher for the note title, using the given fixed/glob/regex patterns.
func Title(f, g, r string) (Matcher, error) {
	return text(f, g, r, func(n *note.Note) string { return n.Frontmatter.Title })
}
