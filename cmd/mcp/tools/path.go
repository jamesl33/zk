package tools

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/jamesl33/zk/internal/vault"
)

// checkPath returns an error unless the current directory is inside a vault and path (relative to
// the current directory) is within that vault, so tools can't reach files outside of it. Symlinks
// are resolved, so a link inside the vault can't be used to reach a file outside of it.
func checkPath(path string) error {
	root, err := vault.RootAbs(".")
	if err != nil {
		return err
	}

	root, err = resolve(root)
	if err != nil {
		return fmt.Errorf("failed to resolve vault root: %w", err)
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %q: %w", path, err)
	}

	abs, err = resolve(abs)
	if err != nil {
		return fmt.Errorf("failed to resolve path %q: %w", path, err)
	}

	if !filepath.IsLocal(mustRel(root, abs)) {
		return fmt.Errorf("path %q is outside the vault", path)
	}

	return nil
}

// resolve returns the given absolute path with symlinks resolved. The path itself needn't exist (e.g. a note which is
// about to be created); in that case its closest existing parent is resolved.
func resolve(abs string) (string, error) {
	resolved, err := filepath.EvalSymlinks(abs)
	if err == nil {
		return resolved, nil
	}

	parent := filepath.Dir(abs)

	if !errors.Is(err, fs.ErrNotExist) || parent == abs {
		return "", err
	}

	resolved, err = resolve(parent)
	if err != nil {
		return "", err
	}

	return filepath.Join(resolved, filepath.Base(abs)), nil
}

// mustRel returns abs relative to root, or an escaping path if that isn't possible.
func mustRel(root, abs string) string {
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return ".."
	}

	return rel
}
