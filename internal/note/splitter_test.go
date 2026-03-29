package note

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitter(t *testing.T) {
	input := "---\ntitle: Test\n---\nBody content"

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(splitter())

	// First scan should get the first marker
	assert.True(t, scanner.Scan())
	assert.Equal(t, "---\n", scanner.Text())

	// Second scan should get the frontmatter AND the second marker
	assert.True(t, scanner.Scan())
	assert.Equal(t, "title: Test\n---\n", scanner.Text())

	// Third scan should get the body
	assert.True(t, scanner.Scan())
	assert.Equal(t, "Body content", scanner.Text())

	assert.False(t, scanner.Scan())
}

func TestSplitterNoFrontmatter(t *testing.T) {
	input := "Just body"

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(splitter())

	// Should just return the whole thing
	assert.True(t, scanner.Scan())
	assert.Equal(t, "Just body", scanner.Text())
	assert.False(t, scanner.Scan())
}
