package matcher

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
	"github.com/jamesl33/zk/internal/note"
)

// Matcher - TODO
type Matcher func(n *note.Note) bool

// NewTitle - TODO
func NewTitle(f, g, r string) (Matcher, error) {
	return matcher(f, g, r, func(n *note.Note) string { return n.Frontmatter.Title })
}

// NewBody - TODO
func NewBody(f, g, r string) (Matcher, error) {
	return matcher(f, g, r, func(n *note.Note) string { return n.Body })
}

// matcher - TODO
func matcher(f, g, r string, extract func(n *note.Note) string) (Matcher, error) {
	matchers := make([]Matcher, 0)

	if m := fixed(f, extract); m != nil {
		matchers = append(matchers, m)
	}

	m, err := glb(g, extract)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	if m != nil {
		matchers = append(matchers, m)
	}

	m, err = regex(r, extract)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	if m != nil {
		matchers = append(matchers, m)
	}

	// TODO
	all := and(matchers...)

	return all, nil
}

// and - TODO
func and(matchers ...Matcher) Matcher {
	return func(n *note.Note) bool {
		for _, matcher := range matchers {
			if !matcher(n) {
				return false
			}
		}

		return true
	}
}

// fixed - TODO
func fixed(pattern string, extract func(n *note.Note) string) Matcher {
	if pattern == "" {
		return nil
	}

	return func(n *note.Note) bool { return strings.Contains(extract(n), pattern) }
}

// glb - TODO
func glb(pattern string, extract func(n *note.Note) string) (Matcher, error) {
	if pattern == "" {
		return nil, nil
	}

	parsed, err := glob.Compile("*" + pattern + "*")
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return func(n *note.Note) bool { return parsed.Match(extract(n)) }, nil
}

// regex - TODO
func regex(pattern string, extract func(n *note.Note) string) (Matcher, error) {
	if pattern == "" {
		return nil, nil
	}

	parsed, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return func(n *note.Note) bool { return parsed.MatchString(extract(n)) }, nil
}
