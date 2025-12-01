package matcher

import (
	"fmt"

	"github.com/jamesl33/zk/internal/note"
)

// text returns a text matcher for the given fixed/glob/regex patterns.
func text(f, g, r string, extract func(n *note.Note) string) (Matcher, error) {
	matchers := make([]Matcher, 0)

	if m := Fixed(f, extract); m != nil {
		matchers = append(matchers, m)
	}

	m, err := Glob(g, extract)
	if err != nil {
		return nil, fmt.Errorf("failed to create glob matcher: %w", err)
	}

	if m != nil {
		matchers = append(matchers, m)
	}

	m, err = Regex(r, extract)
	if err != nil {
		return nil, fmt.Errorf("failed to create regex matcher: %w", err)
	}

	if m != nil {
		matchers = append(matchers, m)
	}

	// Combine the matchers
	all := And(matchers...)

	return all, nil
}
