package searcher

import "strings"

// hidden - TODO
func hidden(p string) bool {
	return strings.HasPrefix(p, ".") || strings.Index(p, "/.") > 0
}
