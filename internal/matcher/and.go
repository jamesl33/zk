package matcher

import "github.com/jamesl33/zk/internal/note"

// And - TODO
func And(matchers ...Matcher) Matcher {
	return func(n *note.Note) bool {
		for _, matcher := range matchers {
			if !matcher(n) {
				return false
			}
		}

		return true
	}
}
