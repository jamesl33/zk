package note

import "time"

// id - TODO
func id() string {
	return time.Now().Format("20060102150405")
}
