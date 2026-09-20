package linter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jamesl33/zk/internal/iterator"
	"github.com/jamesl33/zk/internal/lister"
	"github.com/jamesl33/zk/internal/matcher"
	"github.com/jamesl33/zk/internal/note"
	"github.com/jamesl33/zk/internal/regex"
)

// LintError defines a single linting error.
type LintError struct {
	Path    string
	Message string

	// Line is the 1-based line number the error occurs on, or 0 if the error applies to the
	// whole note (e.g. 'orphan-note', 'duplicate-id') rather than a specific line.
	Line int

	// Column is the 1-based column (character offset within Line) the error occurs on, or 0 if
	// Line is also 0.
	Column int
}

// Linter is a struct that contains the logic for linting notes.
type Linter struct{}

// NewLinter creates a new Linter.
func NewLinter() *Linter {
	return &Linter{}
}

// archiveDir is the directory holding archived notes; it's excluded from linting by default
// since archived notes are considered read-only and their issues aren't worth fixing.
const archiveDir = "4 Archives"

// LintOptions encapsulates the options for linting notes.
type LintOptions struct {
	// archives controls whether notes under archiveDir are included in the results.
	archives bool
}

// WithArchives includes notes under archiveDir in the linting results; they're excluded by
// default.
func WithArchives() func(*LintOptions) {
	return func(o *LintOptions) {
		o.archives = true
	}
}

// Lint performs linting of notes and returns a slice of linting errors.
func (l *Linter) Lint(ctx context.Context, path string, opts ...func(*LintOptions)) ([]*LintError, error) {
	var o LintOptions

	for _, opt := range opts {
		opt(&o)
	}

	ids, paths, errors, err := lintOrphans(ctx, path)
	if err != nil {
		return nil, err
	}

	errors = append(errors, lintDuplicateIDs(paths)...)

	linkErrors, err := lintBrokenLinks(ctx, path, ids)
	if err != nil {
		return nil, err
	}

	errors = append(errors, linkErrors...)

	if !o.archives {
		errors = slices.DeleteFunc(errors, func(e *LintError) bool { return inArchive(e.Path) })
	}

	return errors, nil
}

// inArchive returns a boolean indicating whether path has archiveDir as one of its components.
func inArchive(path string) bool {
	for _, part := range strings.Split(filepath.ToSlash(path), "/") {
		if part == archiveDir {
			return true
		}
	}

	return false
}

// lintOrphans lists every note under path, flagging permanent notes that link to nothing
// (orphan-note). It also returns every note ID and the paths using each ID, so callers don't
// need a second pass over the vault to check for duplicate IDs.
func lintOrphans(ctx context.Context, path string) ([]string, map[string][]string, []*LintError, error) {
	var (
		ids    = make([]string, 0)
		paths  = make(map[string][]string)
		errors = make([]*LintError, 0)
	)

	lstr, err := lister.NewLister(
		lister.WithPath(path),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create lister: %w", err)
	}

	err = iterator.ForEach2(lstr.Many(ctx), func(n *note.Note) error {
		ids = append(ids, n.Name())
		paths[n.Name()] = append(paths[n.Name()], n.Path)

		// A permanent note that links to nothing is a dead end: it can't be reached from, or
		// built on, the rest of the vault.
		if n.Frontmatter.Type == note.Type("permanent") {
			links, err := n.Links()
			if err != nil {
				return fmt.Errorf("failed to get links: %w", err)
			}

			if len(links) == 0 {
				err := LintError{
					Path:    n.Path,
					Message: "Permanent note has no links (orphan-note)",
				}

				errors = append(errors, &err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to list notes: %w", err)
	}

	return ids, paths, errors, nil
}

// lintDuplicateIDs flags every path sharing an ID with another note (duplicate-id).
func lintDuplicateIDs(paths map[string][]string) []*LintError {
	errors := make([]*LintError, 0)

	dupes := make([]string, 0, len(paths))

	for id := range paths {
		dupes = append(dupes, id)
	}

	slices.Sort(dupes)

	for _, id := range dupes {
		ps := paths[id]

		if len(ps) <= 1 {
			continue
		}

		for _, p := range ps {
			err := LintError{
				Path:    p,
				Message: fmt.Sprintf("Identifier %q is used by multiple notes: %s (duplicate-id)", id, strings.Join(ps, ", ")),
			}

			errors = append(errors, &err)
		}
	}

	return errors
}

// lintBrokenLinks lists every note under path, flagging links to IDs not present in ids
// (linkcheck).
func lintBrokenLinks(ctx context.Context, path string, ids []string) ([]*LintError, error) {
	errors := make([]*LintError, 0)

	entire, err := matcher.Entire("", "", regex.Link.String())
	if err != nil {
		return nil, fmt.Errorf("failed to create entire matcher: %w", err)
	}

	lstr, err := lister.NewLister(
		lister.WithPath(path),
		lister.WithMatcher(entire),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create lister: %w", err)
	}

	idx := regex.Link.SubexpIndex("link")

	err = iterator.ForEach2(lstr.Many(ctx), func(n *note.Note) error {
		raw, err := os.ReadFile(n.Path)
		if err != nil {
			return fmt.Errorf("failed to read note: %w", err)
		}

		body := string(raw)

		for _, match := range regex.Link.FindAllStringSubmatchIndex(body, -1) {
			name := body[match[2*idx]:match[2*idx+1]]

			if slices.Contains(ids, name) {
				continue
			}

			line, col := position(body, match[0])

			err := LintError{
				Path:    n.Path,
				Line:    line,
				Column:  col,
				Message: fmt.Sprintf("Link %q is broken (linkcheck)", name),
			}

			errors = append(errors, &err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list notes: %w", err)
	}

	return errors, nil
}

// position returns the 1-based line and column for the given byte offset within body.
func position(body string, offset int) (int, int) {
	before := body[:offset]

	line := strings.Count(before, "\n") + 1

	col := len(before)
	if idx := strings.LastIndex(before, "\n"); idx != -1 {
		col = len(before) - idx - 1
	}

	return line, col + 1
}
