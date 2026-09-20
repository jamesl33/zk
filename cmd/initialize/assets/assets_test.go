package assets

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopySkills(t *testing.T) {
	tmp := t.TempDir()

	dest := filepath.Join(tmp, "skills")

	err := CopySkills(dest)
	require.NoError(t, err)

	entries, err := os.ReadDir(dest)
	require.NoError(t, err)
	assert.NotEmpty(t, entries)

	// Nested files (e.g. under a skill's 'assets' subdirectory) must also be copied.
	data, err := os.ReadFile(filepath.Join(dest, "designer", "assets", "template.md"))
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}

func TestCopySkillsDestNotWritable(t *testing.T) {
	tmp := t.TempDir()

	// Create a file where a directory is expected, forcing MkdirAll to fail.
	blocker := filepath.Join(tmp, "blocker")
	require.NoError(t, os.WriteFile(blocker, []byte("x"), 0o644))

	err := CopySkills(filepath.Join(blocker, "skills"))
	assert.Error(t, err)
}
