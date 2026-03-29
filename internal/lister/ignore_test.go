package lister

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIgnore(t *testing.T) {
	type test struct {
		name     string
		input    string
		expected bool
	}

	tests := []test{
		{
			name:     "regular markdown file",
			input:    "note.md",
			expected: false,
		},
		{
			name:     "GEMINI.md should be ignored",
			input:    "GEMINI.md",
			expected: true,
		},
		{
			name:     "non-markdown file",
			input:    "image.png",
			expected: true,
		},
		{
			name:     "nested markdown file",
			input:    "subdir/note.md",
			expected: false,
		},
		{
			name:     "GEMINI.md suffix only should not match if it's the filename",
			input:    "NOT_GEMINI.md",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ignore(tt.input))
		})
	}
}
