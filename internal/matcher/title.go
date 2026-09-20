package matcher

import "github.com/jamesl33/zk/internal/note"

// Title returns a matcher for the note title, using the given fixed/glob/regex patterns.
func Title(f, g, r string, insensitive bool) (Matcher, error) {
	return text(f, g, r, insensitive, func(n *note.Note) (string, error) { return n.Frontmatter.Title, nil })
}
