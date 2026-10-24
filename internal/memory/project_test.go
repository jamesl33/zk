package memory

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectKeyWorktreesShareKey(t *testing.T) {
	repo := filepath.Join(t.TempDir(), "my repo")

	git(t, t.TempDir(), "init", "-q", repo)
	git(t, repo, "commit", "-q", "--allow-empty", "-m", "init")

	worktree := filepath.Join(t.TempDir(), "elsewhere")

	git(t, repo, "worktree", "add", "-q", "-b", "other", worktree)

	require.Equal(t, "my-repo", ProjectKey(repo))
	require.Equal(t, "my-repo", ProjectKey(worktree))
}

func TestProjectKeyNotRepository(t *testing.T) {
	require.Equal(t, "plain", ProjectKey(filepath.Join(t.TempDir(), "plain")))
}

func TestProjectKeyNeverEmpty(t *testing.T) {
	require.Equal(t, "unknown", ProjectKey("/"))
}
