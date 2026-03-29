package note

import (
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	type test struct {
		name    string
		parents []string
	}

	tests := []test{
		{
			name:    "single parent",
			parents: []string{"notes"},
		},
		{
			name:    "multiple parents",
			parents: []string{"notes", "permanent"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Path(tt.parents...)
			assert.NotEmpty(t, p)

			dir, name := filepath.Split(p)
			assert.Equal(t, filepath.Join(tt.parents...)+string(filepath.Separator), dir)

			// Match timestamp format YYYYMMDDHHMMSS.md
			matched, err := regexp.MatchString(`^\d{14}\.md$`, name)
			require.NoError(t, err)
			assert.True(t, matched)
		})
	}
}

func TestPathNoParents(t *testing.T) {
	p := Path()
	assert.NotEmpty(t, p)

	// Match timestamp format YYYYMMDDHHMMSS.md
	matched, err := regexp.MatchString(`^\d{14}\.md$`, p)
	require.NoError(t, err)
	assert.True(t, matched)
}
