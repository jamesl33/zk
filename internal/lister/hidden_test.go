package lister

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHidden(t *testing.T) {
	type test struct {
		name     string
		input    string
		expected bool
	}

	tests := []test{
		{
			name:     "current directory",
			input:    ".",
			expected: false,
		},
		{
			name:     "hidden file",
			input:    ".hidden",
			expected: true,
		},
		{
			name:     "hidden directory",
			input:    ".git",
			expected: true,
		},
		{
			name:     "regular file",
			input:    "note.md",
			expected: false,
		},
		{
			name:     "regular directory",
			input:    "notes",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, hidden(tt.input))
		})
	}
}
