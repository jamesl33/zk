package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemanticSearchNotesNoVault(t *testing.T) {
	withEmptyDir(t)

	_, _, err := SemanticSearchNotes(t.Context(), nil, &SemanticSearchNotesInput{Query: "hiking"})
	assert.Error(t, err)
}
