package note

import "path/filepath"

// Path returns the path to a new note.
func Path(parents ...string) string {
	return filepath.Join(filepath.Join(parents...), Name())
}
