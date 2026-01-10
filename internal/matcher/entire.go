package matcher

import (
	"github.com/jamesl33/zk/internal/note"
)

// Entire(  )returns a matcher for the entire note, using the given fixed/glob/regex patterns.
func Entire(f, g, r string) (Matcher, error) {
	return text(f, g, r, func(n *note.Note) string { return n.String() })
}
