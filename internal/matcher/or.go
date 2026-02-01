package matcher

import (
	"github.com/jamesl33/zk/internal/note"
)

// Or combines the given matchers using a logic or.
func Or(matchers ...Matcher) Matcher {
	return func(n *note.Note) (bool, error) {
		if len(matchers) == 0 {
			return true, nil
		}

		for _, matcher := range matchers {
			m, err := matcher(n)
			if err != nil {
				return false, err // Purposefully not wrapped
			}

			if m {
				return true, nil
			}
		}

		return false, nil
	}
}
