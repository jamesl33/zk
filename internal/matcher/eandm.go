package matcher

import (
	"github.com/jamesl33/zk/internal/note"
)

// eandm extracts and matches.
func eandm(extract func(n *note.Note) (string, error), match func(text string) bool) func(n *note.Note) (bool, error) {
	return func(n *note.Note) (bool, error) {
		text, err := extract(n)
		if err != nil {
			return false, err // Purposefully not wrapped
		}

		return match(text), nil
	}
}
