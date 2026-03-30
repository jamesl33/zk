package links

import (
	"context"

	"github.com/jamesl33/zk/internal/note"
)

// Replace links within a note, converting the link into the title of the linked note.
func Replace(ctx context.Context, n *note.Note) error {
	err := process(ctx, n, func(link string, title string, ok bool) string {
		if ok {
			return title
		}

		return link
	})

	return err
}
