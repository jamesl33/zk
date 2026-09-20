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

// RootRel returns the vault root as a path relative to the current working directory. It's
// intended for callers which pass the root to a lister, so that resulting note paths are relative
// rather than absolute, while still allowing the vault to be searched from a subdirectory.
func RootRel(path string) (string, error) {
	root, err := Root(path)
	if err != nil {
		return "", err
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(wd, root)
	if err != nil {
		return "", err
	}

	return rel, nil
}
