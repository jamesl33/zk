package vault

import (
	"errors"
	"os"
	"path/filepath"
)

// ErrNotFound is returned by Root when the given path isn't inside a vault.
var ErrNotFound = errors.New("not in a zk vault (no '.zk' directory found)")

// Root walks up from the given path looking for a '.zk' directory, returning the first directory
// found containing one. This allows commands to resolve the vault root even when invoked from a
// subdirectory (e.g. an editor's LSP client using its own working directory).
func Root(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(abs, ".zk")); err == nil {
			return abs, nil
		}

		parent := filepath.Dir(abs)
		if parent == abs {
			return "", ErrNotFound
		}

		abs = parent
	}
}
