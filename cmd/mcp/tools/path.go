package tools

import (
	"fmt"
	"path/filepath"

	"github.com/jamesl33/zk/internal/vault"
)

// checkPath returns an error unless the current directory is inside a vault and path (relative to
// the current directory) is within that vault, so tools can't reach files outside of it.
func checkPath(path string) error {
	root, err := vault.RootAbs(".")
	if err != nil {
		return err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path %q: %w", path, err)
	}

	if !filepath.IsLocal(mustRel(root, abs)) {
		return fmt.Errorf("path %q is outside the vault", path)
	}

	return nil
}

// mustRel returns abs relative to root, or an escaping path if that isn't possible.
func mustRel(root, abs string) string {
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return ".."
	}

	return rel
}
