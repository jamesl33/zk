package matcher

import "github.com/jamesl33/zk/internal/note"

// Any matches anything.
func Any() Matcher {
	return func(_ *note.Note) (bool, error) {
		return true, nil
	}
}
