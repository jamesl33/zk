package note

import "bytes"

// scan is to be used with a 'Scanner' to extract the frontmatter/body from a note.
func scan(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	const marker = "---\n"

	if i := bytes.Index(data, []byte(marker)); i >= 0 {
		return i + len(marker), data[:i+len(marker)], nil
	}

	return len(data), data, nil
}
