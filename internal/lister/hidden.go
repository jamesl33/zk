package lister

import "strings"

// hidden returns a boolean indicating whether the given note is hidden.
func hidden(p string) bool {
	return strings.HasPrefix(p, ".") || strings.Index(p, "/.") > 0
}
