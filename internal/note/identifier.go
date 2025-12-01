package note

import "time"

// id returns a new note name (identifier).
func id() string {
	return time.Now().Format("20060102150405")
}
