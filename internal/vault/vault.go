package vault

import (
	"os"
	"path/filepath"
)

// Root walks up from the given path looking for a '.zk' or '.git' directory, returning the first
// directory found containing one, or "." if no such directory is found. This allows commands to
// resolve the vault root even when invoked from a subdirectory (e.g. an editor's LSP client using
// its own working directory).
func Root(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		for _, marker := range []string{".zk", ".git"} {
			if _, err := os.Stat(filepath.Join(abs, marker)); err == nil {
				return abs, nil
			}
		}

		parent := filepath.Dir(abs)
		if parent == abs {
			return ".", nil
		}

		abs = parent
	}
}
