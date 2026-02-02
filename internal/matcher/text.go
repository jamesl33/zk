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

	// Escapes any special characters
	glob = regexp.QuoteMeta(glob)

	// Convert question marks
	glob = strings.ReplaceAll(glob, "\\?", ".")

	// Convert wildcards
	glob = strings.ReplaceAll(glob, "\\*", ".*")

	// Remove escape sequences for brackets
	glob = br.ReplaceAllString(glob, "[$1]")

	// For non-escaped opening brackets, handle negations
	glob = negbrac(glob)

	return glob
}

// negbrac rewrites the pattern with '[!...]' replaced with '[^...]' where not escaped.
func negbrac(glob string) string {
	var result strings.Builder

	for i := 0; i < len(glob); i++ {
		next := negchar(glob, i)

		// Ignore the '!' character
		if len(next) != 1 {
			i++
		}

		result.WriteString(next)
	}

	return result.String()
}

// negchar converts the '[!' to '[^' where not escaped.
func negchar(glob string, i int) string {
	if glob[i] != '[' {
		return string(glob[i])
	}

	if i-1 >= 0 && glob[i-1] == '\\' {
		return "["
	}

	if i+1 < len(glob) && glob[i+1] != '!' {
		return "["
	}

	return "[^"
}
