package chunker

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const limit = 4096

const prefix = "title: X | text: "

func chunk(body string) []string {
	c := New(func(s string) string { return prefix + s }, limit)

	for _, b := range Blocks(body) {
		c.Add(b)
	}

	return c.Chunks()
}

func TestChunkSingleChunkForSmallNote(t *testing.T) {
	chunks := chunk("Just a short body.")
	require.Len(t, chunks, 1)
	assert.Contains(t, chunks[0], "Just a short body.")
	assert.Contains(t, chunks[0], "title: X")
}

func TestChunkNoBlocks(t *testing.T) {
	assert.Empty(t, chunk(""))
}

func TestChunkPreservesFencedCodeBlocks(t *testing.T) {
	fence := "```go\nfunc main() {\n\tprintln(\"hi\")\n}\n```"
	para := "Paragraph filler text that is reasonably long to consume budget. "
	body := strings.Repeat(para, 60) + "\n\n" + strings.Repeat(para, 30) + "\n\n" + fence + "\n\nMore text after fence."

	chunks := chunk(body)
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

	chunks := chunk(body)
	require.Greater(t, len(chunks), 1)

	for _, c := range chunks {
		assert.LessOrEqual(t, len(c), limit)
	}
}

func TestChunkHardSplitsOnLineBoundaries(t *testing.T) {
	line := strings.Repeat("x", 99)
	body := "```\n" + strings.Repeat(line+"\n", limit/50) + "```"

	chunks := chunk(body)
	require.Greater(t, len(chunks), 1)

	for _, c := range chunks {
		assert.LessOrEqual(t, len(c), limit)

		for l := range strings.SplitSeq(strings.TrimPrefix(c, prefix), "\n") {
			assert.Contains(t, []string{"```", line}, l, "lines must not be split across chunks")
		}
	}
}

func TestChunkLimitCountsCharactersNotBytes(t *testing.T) {
	// 3 bytes per character, so both blocks fit in a chunk of 4096 characters, but not of 4096 bytes.
	para := strings.Repeat("日", 1000)

	chunks := chunk(para + "\n\n" + para)
	require.Len(t, chunks, 1)
	assert.LessOrEqual(t, utf8.RuneCountInString(chunks[0]), limit)
}
