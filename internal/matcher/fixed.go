package matcher

import (
	"strings"

	"github.com/jamesl33/zk/internal/note"
)

// Fixed returns a fixed pattern matcher.
func Fixed(pattern string, extract func(n *note.Note) (string, error)) Matcher {
	if pattern == "" {
		return nil
	}

	fn := func(n *note.Note) (bool, error) {
		c, err := extract(n)
		if err != nil {
			return false, err // Purposefully not wrapped
		}

		return strings.Contains(c, pattern), nil
	}

	return fn
}
