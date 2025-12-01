package matcher

import (
	"fmt"

	"github.com/gobwas/glob"
	"github.com/jamesl33/zk/internal/note"
)

// Glob returns a matcher which will match a given glob pattern.
func Glob(pattern string, extract func(n *note.Note) string) (Matcher, error) {
	if pattern == "" {
		return nil, nil
	}

	parsed, err := glob.Compile("*" + pattern + "*")
	if err != nil {
		return nil, fmt.Errorf("failed to compile glob pattern: %w", err)
	}

	return func(n *note.Note) bool { return parsed.Match(extract(n)) }, nil
}
