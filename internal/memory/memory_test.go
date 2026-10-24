package memory

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/require"
)

// withVault creates a vault and runs in it for the duration of the test.
func withVault(t *testing.T) string {
	t.Helper()

	vault := t.TempDir()

	require.NoError(t, EnsureVault(vault))
	t.Chdir(vault)

	return vault
}

// writeNote writes a permanent note to the vault, returning its path.
func writeNote(t *testing.T, path, title string, tags []string, body string) string {
	t.Helper()

	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))

	content := "---\ntype: permanent\ntitle: " + title + "\ndate: \"2026-01-02\"\ntags: [" + strings.Join(tags, ", ") + "]\n---\n" + body + "\n"

	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	return path
}

// finderOf returns a Finder which finds the notes at the paths, whatever the query, if the threshold admits them.
func finderOf(t *testing.T, distance float32, paths ...string) Finder {
	t.Helper()

	return func(_ context.Context, _ string, limit int, maxDistance float32) ([]*note.Note, error) {
		if distance > maxDistance {
			return nil, nil
		}

		var found []*note.Note

		for _, path := range paths[:min(limit, len(paths))] {
			n, err := note.New(path)
			require.NoError(t, err)

			found = append(found, n)
		}

		return found, nil
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false"}, args...)...)

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
}
