package claude

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jamesl33/zk/cmd/initialize/assets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withTempDir(t *testing.T) string {
	t.Helper()

	tmp := t.TempDir()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmp))
	t.Cleanup(func() { require.NoError(t, os.Chdir(cwd)) })

	return tmp
}

func TestClaudeRun(t *testing.T) {
	tmp := withTempDir(t)

	var (
		c   Claude
		out bytes.Buffer
	)

	require.NoError(t, c.Run(&out))

	data, err := os.ReadFile(filepath.Join(tmp, "ZK.md"))
	require.NoError(t, err)
	assert.Equal(t, assets.Instructions, data)

	data, err = os.ReadFile(filepath.Join(tmp, "CLAUDE.md"))
	require.NoError(t, err)
	assert.Equal(t, "@ZK.md\n", string(data))

	assert.NoDirExists(t, filepath.Join(tmp, ".claude"))
	assert.Contains(t, out.String(), "/plugin marketplace add jamesl33/zk")
}

func TestClaudeRunPrependsImport(t *testing.T) {
	tmp := withTempDir(t)

	require.NoError(t, os.WriteFile("CLAUDE.md", []byte("# Mine\n"), 0o644))

	var c Claude

	require.NoError(t, c.Run(&bytes.Buffer{}))

	data, err := os.ReadFile(filepath.Join(tmp, "CLAUDE.md"))
	require.NoError(t, err)
	assert.Equal(t, "@ZK.md\n\n# Mine\n", string(data))
}

func TestClaudeRunIdempotent(t *testing.T) {
	tmp := withTempDir(t)

	require.NoError(t, os.WriteFile("CLAUDE.md", []byte("@ZK.md\n\n# Mine\n"), 0o644))

	var c Claude

	require.NoError(t, c.Run(&bytes.Buffer{}))
	require.NoError(t, c.Run(&bytes.Buffer{}))

	data, err := os.ReadFile(filepath.Join(tmp, "CLAUDE.md"))
	require.NoError(t, err)
	assert.Equal(t, "@ZK.md\n\n# Mine\n", string(data))
}

func TestClaudeRunMovesImportToTop(t *testing.T) {
	tmp := withTempDir(t)

	require.NoError(t, os.WriteFile("CLAUDE.md", []byte("# Mine\n\n@ZK.md\n\nMore\n"), 0o644))

	var c Claude

	require.NoError(t, c.Run(&bytes.Buffer{}))

	data, err := os.ReadFile(filepath.Join(tmp, "CLAUDE.md"))
	require.NoError(t, err)
	assert.Equal(t, "@ZK.md\n\n# Mine\n\nMore\n", string(data))
}

func TestClaudeRunNewlineAfterImport(t *testing.T) {
	tests := map[string]struct {
		existing string
		expected string
	}{
		"import only, no trailing newline": {existing: "@ZK.md", expected: "@ZK.md\n"},
		"import directly followed by text": {existing: "@ZK.md\n# Mine\n", expected: "@ZK.md\n\n# Mine\n"},
		"text without trailing newline":    {existing: "# Mine", expected: "@ZK.md\n\n# Mine"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tmp := withTempDir(t)

			require.NoError(t, os.WriteFile("CLAUDE.md", []byte(test.existing), 0o644))

			var c Claude

			require.NoError(t, c.Run(&bytes.Buffer{}))

			data, err := os.ReadFile(filepath.Join(tmp, "CLAUDE.md"))
			require.NoError(t, err)
			assert.Equal(t, test.expected, string(data))
		})
	}
}

func TestClaudeRunKeepsClaudeDirectory(t *testing.T) {
	tmp := withTempDir(t)

	skill := filepath.Join(tmp, ".claude", "skills", "mine", "SKILL.md")
	require.NoError(t, os.MkdirAll(filepath.Dir(skill), 0o755))
	require.NoError(t, os.WriteFile(skill, []byte("x"), 0o644))

	var c Claude

	require.NoError(t, c.Run(&bytes.Buffer{}))

	assert.FileExists(t, skill)
}
