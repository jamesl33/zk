package matcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jamesl33/zk/internal/note"
)

// text returns a text matcher for the given fixed/glob/regex patterns.
func text(f, g, r string, extract func(n *note.Note) (string, error)) (Matcher, error) {
	patterns := make([]string, 0, 3)

	if f != "" {
		patterns = append(patterns, regexp.QuoteMeta(f))
	}

	if g != "" {
		patterns = append(patterns, gtor(g))
	}

	if r != "" {
		patterns = append(patterns, r)
	}

	if len(patterns) == 0 {
		return Any(), nil
	}

	// Enable multi-line search
	pattern := fmt.Sprintf(
		"(?m:%s)",
		strings.Join(patterns, "|"),
	)

	parsed, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regular expression: %w", err)
	}

	return eandm(extract, parsed.MatchString), nil
}

// gtor returns a regular expression which is functionally equivalent to the provided glob pattern.
//
// https://en.wikipedia.org/wiki/Glob_(programming)
func gtor(glob string) string {
	// Matches brackets
	br := regexp.MustCompile(`(?U:\\\[(.+)\\\])`)

	// Matches bracket negations
	nbr := regexp.MustCompile(`(?U:((^|[^\\])\[)!)`)

	// Escapes any special characters
	glob = regexp.QuoteMeta(glob)

	// Convert question marks
	glob = strings.ReplaceAll(glob, "\\?", ".")

	// Convert wildcards
	glob = strings.ReplaceAll(glob, "\\*", ".*")

	// Remove escape sequences for brackets
	glob = br.ReplaceAllString(glob, "[$1]")

	// For non-escaped opening brackets, handle negations
	glob = nbr.ReplaceAllString(glob, "$1^")

	return glob
}
