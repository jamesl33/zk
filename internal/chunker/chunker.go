// Package chunker splits notes into chunks which fit within a character limit.
package chunker

import (
	"strings"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Chunker accumulates markdown blocks into chunks that each fit under a character limit. The
// frontmatter is prepended to every chunk.
type Chunker struct {
	frontmatter string
	budget      int
	chunks      []string
	current     strings.Builder
}

// New returns a chunker which produces chunks of at most limit characters.
func New(frontmatter string, limit int) *Chunker {
	budget := limit - len(frontmatter) - len("\n\n")
	if budget <= 0 {
		budget = limit
	}

	chunker := Chunker{
		frontmatter: frontmatter,
		budget:      budget,
	}

	return &chunker
}

// Add appends the given block to the current chunk, starting a new chunk if it wouldn't fit.
func (c *Chunker) Add(b string) {
	if len(b) > c.budget {
		c.flush()

		for _, piece := range hardSplit(b, c.budget) {
			c.chunks = append(c.chunks, c.frontmatter+"\n\n"+piece)
		}

		return
	}

	if c.current.Len() > 0 && c.current.Len()+len("\n\n")+len(b) > c.budget {
		c.flush()
	}

	if c.current.Len() > 0 {
		c.current.WriteString("\n\n")
	}

	c.current.WriteString(b)
}

// Chunks completes the current chunk and returns all the chunks.
func (c *Chunker) Chunks() []string {
	c.flush()

	// No body content (e.g. a frontmatter-only note with tags but an empty body) -- embed the
	// frontmatter alone.
	if len(c.chunks) == 0 {
		return []string{c.frontmatter}
	}

	return c.chunks
}

// flush completes the current chunk, if any.
func (c *Chunker) flush() {
	if c.current.Len() == 0 {
		return
	}

	c.chunks = append(c.chunks, c.frontmatter+"\n\n"+c.current.String())
	c.current.Reset()
}

// Blocks returns the raw source of each top-level markdown block in body.
//
// Blocks are sliced between consecutive Pos() offsets, not rebuilt from Lines(), which omits
// delimiters such as code fences.
func Blocks(body string) []string {
	if strings.TrimSpace(body) == "" {
		return nil
	}

	var (
		src    = []byte(body)
		doc    = goldmark.DefaultParser().Parse(text.NewReader(src))
		blocks []string
	)

	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		if block := blockSource(src, n); block != "" {
			blocks = append(blocks, block)
		}
	}

	return blocks
}

// blockSource returns the raw source of the given top-level block, running up to the start of the
// next block (or the end of src).
func blockSource(src []byte, n ast.Node) string {
	start := n.Pos()
	if start < 0 {
		return ""
	}

	stop := len(src)
	if next := n.NextSibling(); next != nil && next.Pos() >= 0 {
		stop = next.Pos()
	}

	return strings.TrimRight(string(src[start:stop]), "\n")
}

// hardSplit splits a single oversized block into pieces at most budget runes long, as a fallback for the rare case a
// block alone exceeds budget (e.g. a very large code fence).
//
// Pieces are cut at the last line break within budget, where possible, to avoid splitting a word or identifier across
// two pieces.
func hardSplit(s string, budget int) []string {
	if budget <= 0 {
		budget = 1
	}

	var (
		runes  = []rune(s)
		chunks = make([]string, 0, len(runes)/budget+1)
	)

	for len(runes) > budget {
		var (
			n, skip = budget, 0
			prefix  = string(runes[:budget+1])
		)

		if i := strings.LastIndexByte(prefix, '\n'); i > 0 {
			n, skip = utf8.RuneCountInString(prefix[:i]), 1
		}

		chunks = append(chunks, string(runes[:n]))

		runes = runes[n+skip:]
	}

	if len(runes) > 0 {
		chunks = append(chunks, string(runes))
	}

	return chunks
}
