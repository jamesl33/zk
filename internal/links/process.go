package links

import (
	"context"
	"fmt"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/regex"
)

// process runs the given function to perform replacements for all the links within the given note.
func process(
	ctx context.Context,
	n *note.Note,
	fn func(link string, title string, ok bool) string,
) error {
	body, err := n.GetBody()
	if err != nil {
		return fmt.Errorf("failed to get body: %w", err)
	}

	links, err := n.Links()
	if err != nil {
		return fmt.Errorf("failed to get links: %w", err)
	}

	if len(links) == 0 {
		return nil
	}

	matchers := hs.Map(links, func(n string) matcher.Matcher { return matcher.Name(n) })

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return fmt.Errorf("failed to create lister: %w", err)
	}

	titles := make(map[string]string)

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		titles[n.Name()] = n.Frontmatter.Title
	}))
	if err != nil {
		return fmt.Errorf("failed to list notes: %w", err)
	}

	// rpl replaces links within notes.
	rpl := func(match string) string {
		var (
			submatches = regex.Link.FindStringSubmatch(match)
			link       = submatches[regex.Link.SubexpIndex("link")]
		)

		title, ok := titles[link]

		return fn(link, title, ok)
	}

	// Rewrite the note body
	n.SetBody(regex.Link.ReplaceAllStringFunc(body, rpl))

	return nil
}
