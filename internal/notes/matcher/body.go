package matcher

import "github.com/jamesl33/zk/internal/note"

// Body - TODO
func Body(f, g, r string) (Matcher, error) {
	return text(f, g, r, func(n *note.Note) string { return n.Body })
}
