package notes

import "strings"

// hidden - TODO
func hidden(p string) bool {
	return strings.Index(p, "/.") > 0
}
