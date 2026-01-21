package links

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/ptr"
	"github.com/jamesl33/zk/internal/regex"
)

// Replace links within a note, converting the link into the title of the linked note.
//
// NOTE: Returns a shallow copy of the given note, with a re-written `Body`.
func Replace(ctx context.Context, a *note.Note) (*note.Note, error) {
	var (
		b     = ptr.To(*a)
		links = b.Links()
	)

	if len(links) == 0 {
		return b, nil
	}

	matchers := hs.Map(links, func(n string) matcher.Matcher { return matcher.Name(n) })

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	titles := make(map[string]string)

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		titles[n.Name()] = n.Frontmatter.Title
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	// rpl replaces links within notes, with the target notes title.
	rpl := func(match string) string {
		var (
			submatches = regex.Link.FindStringSubmatch(match)
			link       = submatches[regex.Link.SubexpIndex("link")]
		)

		if title, ok := titles[link]; ok {
			return title
		}

		return link
	}

	// Rewrite the note body
	b.Body = regex.Link.ReplaceAllStringFunc(b.Body, rpl)

	return b, nil
}
