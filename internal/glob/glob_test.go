package glob

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToRegexpMatch(t *testing.T) {
	type test struct {
		glob  string
		input string
	}

	tests := []test{
		{glob: "*.md", input: "note.md"},
		{glob: "note?.md", input: "note1.md"},
		{glob: "note[12].md", input: "note1.md"},
		{glob: "note[12].md", input: "note2.md"},
		{glob: "note[!1].md", input: "note2.md"},
		{glob: "note[?].md", input: "note?.md"},
		{glob: "note[*].md", input: "note*.md"},
		{glob: "note\\*.md", input: "note*.md"},
		{glob: "note\\?.md", input: "note?.md"},
		{glob: "note\\[1\\].md", input: "note[1].md"},
	}

	for _, tt := range tests {
		t.Run(tt.glob+"_"+tt.input, func(t *testing.T) {
			pattern := ToRegexp(tt.glob)
			assert.Regexp(t, "^"+pattern+"$", tt.input)
		})
	}
}

func TestToRegexpNoMatch(t *testing.T) {
	type test struct {
		glob  string
		input string
	}

	tests := []test{
		{glob: "*.md", input: "note.txt"},
		{glob: "note?.md", input: "note.md"},
		{glob: "note[12].md", input: "note3.md"},
		{glob: "note[!1].md", input: "note1.md"},
		{glob: "note[?].md", input: "note1.md"},
		{glob: "note\\*.md", input: "notex.md"},
	}

	for _, tt := range tests {
		t.Run(tt.glob+"_"+tt.input, func(t *testing.T) {
			pattern := ToRegexp(tt.glob)
			assert.NotRegexp(t, "^"+pattern+"$", tt.input)
		})
	}
}
