package matcher

import (
	"github.com/jamesl33/zk/internal/note"
)

// Not negates a matcher.
func Not(matcher Matcher) Matcher {
	fn := func(n *note.Note) (bool, error) {
		m, err := matcher(n)
		if err != nil {
			return false, err // Purposefully not wrapped
		}

		return !m, nil
	}

	return fn
}
