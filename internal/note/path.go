package note

import "path/filepath"

// Path - TODO
func Path(parents ...string) string {
	return filepath.Join(filepath.Join(parents...), Name())
}
