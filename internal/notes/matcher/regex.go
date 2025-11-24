package matcher

import (
	"fmt"
	"regexp"

	"github.com/jamesl33/zk/internal/note"
)

// Regex - TODO
func Regex(pattern string, extract func(n *note.Note) string) (Matcher, error) {
	if pattern == "" {
		return nil, nil
	}

	parsed, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return func(n *note.Note) bool { return parsed.MatchString(extract(n)) }, nil
}
