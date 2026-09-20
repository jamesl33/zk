// Package glob converts shell-style glob patterns into equivalent regular expressions.
package glob

import (
	"regexp"
	"strings"
)

// ToRegexp returns a regular expression which is functionally equivalent to the provided glob pattern.
//
// https://en.wikipedia.org/wiki/Glob_(programming)
func ToRegexp(glob string) string {
	runes := []rune(glob)

	var out strings.Builder

	for i := 0; i < len(runes); i++ {
		switch c := runes[i]; c {
		case '\\':
			i = writeEscape(&out, runes, i)
		case '?':
			out.WriteString(".")
		case '*':
			out.WriteString(".*")
		case '[':
			if j, ok := writeBracket(&out, runes, i); ok {
				i = j
				continue
			}

			out.WriteString(regexp.QuoteMeta(string(c)))
		default:
			out.WriteString(regexp.QuoteMeta(string(c)))
		}
	}

	return out.String()
}

// writeEscape writes the literal character following a backslash at runes[i], and returns the index
// of the last rune it consumed.
func writeEscape(out *strings.Builder, runes []rune, i int) int {
	if i+1 >= len(runes) {
		out.WriteString(regexp.QuoteMeta(string(runes[i])))
		return i
	}

	out.WriteString(regexp.QuoteMeta(string(runes[i+1])))

	return i + 1
}

// writeBracket writes the regular expression bracket expression starting at the '[' at runes[i]. It
// returns the index of the closing ']' and true, or false if runes[i:] has no closing bracket.
func writeBracket(out *strings.Builder, runes []rune, i int) (int, bool) {
	j := i + 1

	negate := j < len(runes) && runes[j] == '!'
	if negate {
		j++
	}

	start := j
	for j < len(runes) && runes[j] != ']' {
		j++
	}

	if j >= len(runes) {
		return i, false
	}

	out.WriteString("[")
	if negate {
		out.WriteString("^")
	}

	writeBracketContent(out, runes[start:j])

	out.WriteString("]")

	return j, true
}

// writeBracketContent writes the contents of a bracket expression literally, escaping any characters
// which would otherwise be given special meaning by the regular expression engine.
func writeBracketContent(out *strings.Builder, content []rune) {
	for _, c := range content {
		if c == '\\' || c == '^' {
			out.WriteString("\\")
		}

		out.WriteRune(c)
	}
}
