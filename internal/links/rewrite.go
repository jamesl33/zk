package links

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/note"
)

// Rewrite the links in a note so that links contain up-to-date titles.
func Rewrite(ctx context.Context, n *note.Note) error {
	err := process(ctx, n, func(link string, title string, ok bool) string {
		if ok {
			return fmt.Sprintf("[[%s|%s]]", link, title)
		}

		return fmt.Sprintf("[[%s]]", link)
	})

	return err
}
