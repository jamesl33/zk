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

// Rewrite the links in a note so that links contain up-to-date titles.
//
// TODO (jamesl33): There's a huge amount of duplication between rewrite and replace.
func Rewrite(ctx context.Context, n *note.Note) error {
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

	// rpl replaces links within notes, with the target notes title.
	rpl := func(match string) string {
		var (
			submatches = regex.Link.FindStringSubmatch(match)
			link       = submatches[regex.Link.SubexpIndex("link")]
		)

		if title, ok := titles[link]; ok {
			return fmt.Sprintf("[[%s|%s]]", link, title)
		}

		return fmt.Sprintf("[[%s]]", link)
	}

	// Rewrite the note body
	n.SetBody(regex.Link.ReplaceAllStringFunc(body, rpl))

	return nil
}
