package note

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestName(t *testing.T) {
	n := Name()
	assert.NotEmpty(t, n)
	assert.Contains(t, n, ".md")

	// Match timestamp format YYYYMMDDHHMMSS.md
	matched, err := regexp.MatchString(`^\d{14}\.md$`, n)
	require.NoError(t, err)
	assert.True(t, matched)
}
