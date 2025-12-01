package matcher

import (
	"strings"

	"github.com/jamesl33/zk/internal/note"
)

// Fixed returns a fixed pattern matcher.
func Fixed(pattern string, extract func(n *note.Note) string) Matcher {
	if pattern == "" {
		return nil
	}

	return func(n *note.Note) bool { return strings.Contains(extract(n), pattern) }
}
