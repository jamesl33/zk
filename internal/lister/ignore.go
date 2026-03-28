package lister

import (
	"path/filepath"
	"strings"
)

// ignore returns a boolean indicating whether the given file should be skipped.
func ignore(p string) bool {
	if !strings.HasSuffix(p, ".md") {
		return true
	}

	switch filepath.Base(p) {
	case "GEMINI.md":
		return true
	}

	return false
}
