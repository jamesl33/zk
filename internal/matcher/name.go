package matcher

import "github.com/jamesl33/zk/internal/note"

// Name returns a matcher for the given note name.
func Name(name string) Matcher {
	return func(n *note.Note) bool { return name == n.Name() }
}
