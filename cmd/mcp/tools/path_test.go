package tools

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
