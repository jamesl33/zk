package matcher

import "github.com/jamesl33/zk/internal/note"

// Name - TODO
func Name(name string) Matcher {
	return func(n *note.Note) bool { return name == n.Name() }
}
