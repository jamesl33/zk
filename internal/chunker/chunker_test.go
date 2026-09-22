package chunker

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const limit = 4096

func chunk(frontmatter, body string) []string {
	c := New(frontmatter, limit)

	for _, b := range Blocks(body) {
		c.Add(b)
	}

	return c.Chunks()
}

func TestChunkSingleChunkForSmallNote(t *testing.T) {
	chunks := chunk("---\ntitle: X\n---", "Just a short body.")
	require.Len(t, chunks, 1)
	assert.Contains(t, chunks[0], "Just a short body.")
	assert.Contains(t, chunks[0], "title: X")
}

func TestChunkFrontmatterOnlyNote(t *testing.T) {
	chunks := chunk("---\ntags: [a]\n---", "")
	assert.Equal(t, []string{"---\ntags: [a]\n---"}, chunks)
}

func TestChunkPreservesFencedCodeBlocks(t *testing.T) {
	fence := "```go\nfunc main() {\n\tprintln(\"hi\")\n}\n```"
	para := "Paragraph filler text that is reasonably long to consume budget. "
	body := strings.Repeat(para, 60) + "\n\n" + strings.Repeat(para, 30) + "\n\n" + fence + "\n\nMore text after fence."

	chunks := chunk("---\ntitle: X\n---", body)
	require.Greater(t, len(chunks), 1)

	found := false

	for _, c := range chunks {
		if strings.Contains(c, fence) {
			found = true
		}
	}

	assert.True(t, found, "fenced code block must not be split across chunks")
}

func TestChunkHardSplitsOversizedBlock(t *testing.T) {
	body := strings.Repeat("x", limit*2)

	chunks := chunk("---\ntitle: X\n---", body)
	require.Greater(t, len(chunks), 1)

	for _, c := range chunks {
		assert.LessOrEqual(t, len(c), limit)
	}
}
