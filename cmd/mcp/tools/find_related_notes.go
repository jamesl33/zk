package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"

	"github.com/jamesl33/zk/internal/hs"
	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/vector"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// FindRelatedNotesInput - TODO
type FindRelatedNotesInput struct {
	// Path - TODO
	Path string `json:"path" jsonschema:"The path to a note"`
}

// FindRelatedNotesOutput - TODO
type FindRelatedNotesOutput struct {
	// Notes - TODO
	Notes []*note.Note `json:"notes" jsonschema:"Notes which are semantically similar to the given note"`
}

// FindRelatedNotes - TODO
//
// TODO (jamesl33): De-duplicate this code?
func FindRelatedNotes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input *FindRelatedNotesInput,
) (*mcp.CallToolResult, *FindRelatedNotesOutput, error) {
	n, err := note.New(input.Path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open note: %w", err)
	}

	db, err := vector.New(ctx, filepath.Join(".zk", "zk.sqlite3"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	err = populate(ctx, db)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to populate database: %w", err)
	}

	notes, err := db.Find(ctx, n)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find related notes: %w", err)
	}

	found := make([]*note.Note, 0)

	err = iterator.ForEach2(notes, hs.Infallible(func(n *note.Note) {
		// TODO
		n.Body = ""

		// TODO
		found = append(found, n)
	}))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list related notes: %w", err)
	}

	to, err := to(ctx, n)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err) // TODO
	}

	from, err := from(ctx, n)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err) // TODO
	}

	output := FindRelatedNotesOutput{
		Notes: slices.Concat(found, to, from),
	}

	return nil, &output, nil
}

// to - TODO
func to(ctx context.Context, n *note.Note) ([]*note.Note, error) {
	var (
		// name of the note, escaped for use in regular expressions
		name = regexp.QuoteMeta(n.Name())

		// pattern which matches links to this note
		pattern = fmt.Sprintf(`\[\[%s(\|.*?)?\]\]`, name)
	)

	matcher, err := matcher.Body("", "", pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to create matcher: %w", err)
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	found := make([]*note.Note, 0)

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		// TODO
		n.Body = ""

		// TODO
		found = append(found, n)
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	return found, nil
}

// from - TODO
func from(ctx context.Context, n *note.Note) ([]*note.Note, error) {
	matchers := hs.Map(n.Links(), func(n string) matcher.Matcher { return matcher.Name(n) })

	// Must check for no matchers, as the default is to list all
	if len(matchers) == 0 {
		return nil, nil
	}

	lister, err := lister.NewLister(
		lister.WithPath("."),
		lister.WithMatcher(matcher.Or(matchers...)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	found := make([]*note.Note, 0)

	err = iterator.ForEach2(lister.Many(ctx), hs.Infallible(func(n *note.Note) {
		// TODO
		n.Body = ""

		// TODO
		found = append(found, n)
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	return found, nil
}
