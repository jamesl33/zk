package tools

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckPath(t *testing.T) {
	tmp := withVault(t)

	assert.NoError(t, checkPath("."))
	assert.NoError(t, checkPath("0 Inbox/note.md"))
	assert.NoError(t, checkPath(filepath.Join(tmp, "note.md")))
	assert.Error(t, checkPath(".."))
	assert.Error(t, checkPath("../other/note.md"))
	assert.Error(t, checkPath("/etc/passwd"))
}

func TestCheckPathNoVault(t *testing.T) {
	withEmptyDir(t)

	assert.Error(t, checkPath("."))
}

func TestCheckPathSymlinkEscape(t *testing.T) {
	tmp := withVault(t)

	outside := t.TempDir()

	require.NoError(t, os.Symlink(outside, filepath.Join(tmp, "escape")))
	require.NoError(t, os.WriteFile(filepath.Join(outside, "note.md"), nil, 0o644))
	require.NoError(t, os.Symlink(filepath.Join(outside, "note.md"), filepath.Join(tmp, "link.md")))

	assert.Error(t, checkPath("escape"))
	assert.Error(t, checkPath("escape/note.md"))
	assert.Error(t, checkPath("escape/new/note.md"))
	assert.Error(t, checkPath("link.md"))
}

func TestCheckPathNotYetCreated(t *testing.T) {
	withVault(t)

	assert.NoError(t, checkPath("new dir/deeper/note.md"))
}
