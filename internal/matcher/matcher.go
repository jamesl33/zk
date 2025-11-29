package matcher

import (
	"github.com/jamesl33/zk/internal/note"
)

// Matcher - TODO
type Matcher func(n *note.Note) bool
