package matcher

import "github.com/jamesl33/zk/internal/note"

// Body returns a matcher for the note body, using the given fixed/glob/regex patterns.
func Body(f, g, r string, insensitive bool) (Matcher, error) {
	return text(f, g, r, insensitive, func(n *note.Note) (string, error) { return n.GetBody() })
}
