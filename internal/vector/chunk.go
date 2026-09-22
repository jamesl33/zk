package vector

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

// ollamaEmbedContextTokens is the context window (in tokens) of the default Ollama embedding
// model ("embeddinggemma"). Ollama's Embed API is called with Truncate: false (see
// internal/ai/ollama.go), so oversized input errors instead of being silently truncated -- notes
// that don't fit must be split into multiple chunks and embedded/stored separately.
const ollamaEmbedContextTokens = 2048

// maxChunkChars is the character budget for a single chunk. There's no tokenizer available in
// this codebase to count tokens exactly, so a conservative chars-per-token estimate is used
// instead: plain English prose averages ~4 chars/token, but markdown adds extra short/punctuation
// tokens (#, -, **, [](), etc.) that push the real average down, so 3 chars/token is assumed to
// avoid under-counting tokens. An extra ~20% safety margin is then taken off the resulting budget
// (2048 * 3 = 6144 -> 4800) to leave headroom for tokenizer variance.
const maxChunkChars = 4800

// chunk splits a note's frontmatter and body into one or more strings that each fit under
// maxChunkChars, so long notes can be embedded as multiple vectors instead of erroring against
// Ollama. frontmatter is prepended to every chunk so each embedded piece keeps its title/tags/date
// context. body is split along markdown block boundaries (headings, paragraphs, code fences,
// lists) rather than a raw character/paragraph split, so a chunk boundary never lands inside a
// code fence or list item unless a single block alone exceeds the budget.
func chunk(frontmatter, body string) []string {
	budget := maxChunkChars - len(frontmatter) - len("\n\n")
	if budget <= 0 {
		budget = maxChunkChars
	}

	blocks := markdownBlocks(body)

	var (
		chunks  []string
		current strings.Builder
	)

	flush := func() {
		if current.Len() == 0 {
			return
		}

		chunks = append(chunks, frontmatter+"\n\n"+current.String())
		current.Reset()
	}

	for _, b := range blocks {
		if len(b) > budget {
			flush()

			for _, piece := range hardSplit(b, budget) {
				chunks = append(chunks, frontmatter+"\n\n"+piece)
			}

			continue
		}

		if current.Len() > 0 && current.Len()+len("\n\n")+len(b) > budget {
			flush()
		}

		if current.Len() > 0 {
			current.WriteString("\n\n")
		}

		current.WriteString(b)
	}

	flush()

	// No body content (e.g. a frontmatter-only note with tags but an empty body) -- embed the
	// frontmatter alone, matching the previous single-embed-call behaviour for that case.
	if len(chunks) == 0 {
		chunks = []string{frontmatter}
	}

	return chunks
}

// markdownBlocks parses body as markdown and returns the raw source text of each top-level block
// node (headings, paragraphs, lists, code fences, blockquotes, ...), in document order. Each
// block's span runs from its own start position to the next block's start position (or the end of
// source for the last block) rather than from its Lines() segments: several node kinds --
// notably FencedCodeBlock -- only record their inner content lines in Lines(), excluding the
// delimiter lines (``` fences, list markers, etc.), so reconstructing from Lines() would silently
// drop them. Every block-opening parser sets Pos() to the exact source offset where the block
// starts (see goldmark/parser.Parser.parseBlocks), which does include those delimiters.
func markdownBlocks(body string) []string {
	if strings.TrimSpace(body) == "" {
		return nil
	}

	src := []byte(body)
	doc := goldmark.DefaultParser().Parse(text.NewReader(src))

	var blocks []string

	for n := doc.FirstChild(); n != nil; n = n.NextSibling() {
		start := n.Pos()
		if start < 0 {
			continue
		}

		stop := len(src)
		if next := n.NextSibling(); next != nil && next.Pos() >= 0 {
			stop = next.Pos()
		}

		block := strings.TrimRight(string(src[start:stop]), "\n")
		if block != "" {
			blocks = append(blocks, block)
		}
	}

	return blocks
}

// hardSplit splits a single oversized block into pieces at most budget runes long, as a fallback
// for the rare case a block alone exceeds budget (e.g. a very large code fence).
func hardSplit(s string, budget int) []string {
	if budget <= 0 {
		budget = 1
	}

	runes := []rune(s)
	chunks := make([]string, 0, len(runes)/budget+1)

	for len(runes) > 0 {
		n := min(len(runes), budget)
		chunks = append(chunks, string(runes[:n]))
		runes = runes[n:]
	}

	return chunks
}
