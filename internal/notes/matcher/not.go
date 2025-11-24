package matcher

import "github.com/jamesl33/zk/internal/note"

// Not - TODO
func Not(m Matcher) Matcher {
	return func(n *note.Note) bool { return !m(n) }
}
