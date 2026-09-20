package matcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jamesl33/zk/internal/glob"
	"github.com/jamesl33/zk/internal/note"
)

// text returns a text matcher for the given fixed/glob/regex patterns.
func text(f, g, r string, extract func(n *note.Note) (string, error)) (Matcher, error) {
	patterns := make([]string, 0, 3)

	if f != "" {
		patterns = append(patterns, regexp.QuoteMeta(f))
	}

	if g != "" {
		patterns = append(patterns, glob.ToRegexp(g))
	}

	if r != "" {
		patterns = append(patterns, r)
	}

	if len(patterns) == 0 {
		return Any(), nil
	}

	// Enable multi-line search, and case-insensitive search unless a pattern has an uppercase letter
	flags := "m"
	if insensitive(f, g, r) {
		flags += "i"
	}

	pattern := fmt.Sprintf(
		"(?%s:%s)",
		flags,
		strings.Join(patterns, "|"),
	)

	parsed, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regular expression: %w", err)
	}

	return eandm(extract, parsed.MatchString), nil
}

// insensitive reports whether the given patterns should be matched case-insensitively: insensitive unless one of the
// patterns contains an uppercase letter.
func insensitive(patterns ...string) bool {
	for _, p := range patterns {
		if strings.ToLower(p) != p {
			return false
		}
	}

	return true
}
