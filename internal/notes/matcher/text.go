package matcher

import (
	"fmt"

	"github.com/jamesl33/zk/internal/note"
)

// text - TODO
func text(f, g, r string, extract func(n *note.Note) string) (Matcher, error) {
	matchers := make([]Matcher, 0)

	if m := Fixed(f, extract); m != nil {
		matchers = append(matchers, m)
	}

	m, err := Glob(g, extract)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	if m != nil {
		matchers = append(matchers, m)
	}

	m, err = Regex(r, extract)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	if m != nil {
		matchers = append(matchers, m)
	}

	// TODO
	all := And(matchers...)

	return all, nil
}
