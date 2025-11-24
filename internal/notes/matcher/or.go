package matcher

import "github.com/jamesl33/zk/internal/note"

// Or - TODO
func Or(matchers ...Matcher) Matcher {
	return func(n *note.Note) bool {
		if len(matchers) == 0 {
			return true
		}

		for _, matcher := range matchers {
			if matcher(n) {
				return true
			}
		}

		return false
	}
}
