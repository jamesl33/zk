package matcher

import (
	"github.com/jamesl33/zk/internal/note"
)

// And combines the given matchers using a logic and.
func And(matchers ...Matcher) Matcher {
	return func(n *note.Note) (bool, error) {
		for _, matcher := range matchers {
			m, err := matcher(n)
			if err != nil {
				return false, err // Purposefully not wrapped
			}

			if !m {
				return false, nil
			}
		}

		return true, nil
	}
}
