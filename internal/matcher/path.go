package matcher

import "github.com/jamesl33/zk/internal/note"

// Path returns a matcher for the note path, using the given fixed/glob/regex patterns.
func Path(f, g, r string) (Matcher, error) {
	return text(f, g, r, func(n *note.Note) string { return n.Path })
}
