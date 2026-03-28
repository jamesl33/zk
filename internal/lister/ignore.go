package lister

import (
	"strings"
)

// ignore returns a boolean indicating whether the given file should be skipped.
func ignore(n string) bool {
	if !strings.HasSuffix(n, ".md") {
		return true
	}

	switch n {
	case "GEMINI.md":
		return true
	}

	return false
}
