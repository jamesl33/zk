package lister

import "strings"

// hidden returns a boolean indicating whether the given note is hidden.
func hidden(n string) bool {
	return n != "." && strings.HasPrefix(n, ".")
}
