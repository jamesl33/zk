package matcher

import "github.com/jamesl33/zk/internal/note"

// Path returns a matcher for the note path, using the given fixed/glob/regex patterns.
func Path(f, g, r string, insensitive bool) (Matcher, error) {
	return text(f, g, r, insensitive, func(n *note.Note) (string, error) { return n.Path, nil })
}
