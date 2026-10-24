package memory

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var unsafeKey = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// ProjectKey names the project the directory belongs to: the name of its git repository, shared by every worktree, else
// the directory's own name.
func ProjectKey(cwd string) string {
	name := filepath.Base(cwd)

	out, err := exec.Command("git", "-C", cwd, "rev-parse", "--path-format=absolute", "--git-common-dir").Output()
	if common := strings.TrimSpace(string(out)); err == nil && common != "" {
		name = filepath.Base(common)

		// A normal repository's common dir is '<repo>/.git', a bare one is '<repo>.git'.
		if name == ".git" {
			name = filepath.Base(filepath.Dir(common))
		}

		name = strings.TrimSuffix(name, ".git")
	}

	name = strings.Trim(unsafeKey.ReplaceAllString(name, "-"), "-.")
	if name == "" {
		return "unknown"
	}

	return name
}

// ProjectDir is the directory in the vault holding the project's notes.
func ProjectDir(key string) string {
	return filepath.Join("1 Projects", key)
}
