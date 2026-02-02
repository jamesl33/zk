package note

import "bytes"

// splitter returns a splitter to be used with a 'Scanner' to extract the frontmatter/body from a note.
func splitter() func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	var markers int

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		const marker = "---\n"

		i := bytes.Index(data, []byte(marker))

		if i < 0 || markers >= 2 {
			return len(data), data, nil
		}

		markers++

		return i + len(marker), data[:i+len(marker)], nil
	}
}
