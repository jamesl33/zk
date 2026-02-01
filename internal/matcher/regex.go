package matcher

import (
	"fmt"
	"regexp"

	"github.com/jamesl33/zk/internal/note"
)

// Regex returns a matcher which will match a given regular expression pattern.
func Regex(pattern string, extract func(n *note.Note) (string, error)) (Matcher, error) {
	if pattern == "" {
		return nil, nil
	}

	// Enable multi-line search
	pattern = fmt.Sprintf(
		"(?m:%s)",
		pattern,
	)

	parsed, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regular expression: %w", err)
	}

	return eandm(extract, parsed.MatchString), nil
}
