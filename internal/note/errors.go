package note

import "errors"

// ErrNotNote is returned if a given file doesn't appear to be a note.
var ErrNotNote = errors.New("not a note")
