package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindRelatedNotesNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := FindRelatedNotes(t.Context(), nil, &FindRelatedNotesInput{Path: "note.md"})
	assert.Error(t, err)
}

func TestFindRelatedNotesNotFound(t *testing.T) {
	withVault(t)

	_, _, err := FindRelatedNotes(t.Context(), nil, &FindRelatedNotesInput{Path: "does-not-exist.md"})
	assert.Error(t, err)
}
