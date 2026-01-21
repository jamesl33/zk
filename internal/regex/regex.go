package regex

import "regexp"

// Link matches WikiLink style links between notes in the Zettelkasten.
var Link = regexp.MustCompile(`\[\[(?P<link>\d{14}?)(\|(?P<text>.*?))?\]\]`)
