package matcher

import (
	"github.com/jamesl33/zk/internal/note"
)

// Matcher is a readability wrapper for a function which "matches" a given note (e.g. determines if to select or not).
type Matcher func(n *note.Note) bool
