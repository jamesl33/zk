package regex

import "regexp"

// Link matches WikiLink style links between notes in the Zettelkasten.
var Link = regexp.MustCompile(`\[\[(?P<link>\d{14}?)(\|(?P<text>.*?))?\]\]`)

// MarkdownNoteLink matches markdown style links to a note, which should be WikiLinks instead; the
// target ends in a note ID, optionally followed by '.md'.
var MarkdownNoteLink = regexp.MustCompile(`\[[^\]\n]*\]\((?P<target>[^)\n]*?\d{14}(?:\.md)?)\)`)
