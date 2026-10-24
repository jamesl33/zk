package memory

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/jamesl33/zk/internal/lister"
)

// digestSectionMax is the most bytes of note titles in each section of the digest.
const digestSectionMax = 1000

// globalTags are the tags of notes which are always listed, whichever project they were saved from.
var globalTags = []string{"user", "feedback"}

// Digest is what Claude is told at the start of a session (and again after compaction): how to save memories, and the
// titles of those worth knowing before the first prompt. It runs in the vault.
func Digest(ctx context.Context, key string) (string, error) {
	lister, err := lister.NewLister(lister.WithPath("."))
	if err != nil {
		return "", fmt.Errorf("failed to create lister: %w", err)
	}

	var project, global []listing

	for n, err := range lister.Many(ctx) {
		if err != nil {
			return "", fmt.Errorf("failed to list notes: %w", err)
		}

		l := listing{date: n.Frontmatter.Date, line: fmt.Sprintf("- %s (%s)", n.Frontmatter.Title, n.Path)}

		switch {
		case inProject(n.Path, key):
			project = append(project, l)
		case slices.ContainsFunc(n.Frontmatter.Tags, func(tag string) bool { return slices.Contains(globalTags, tag) }):
			global = append(global, l)
		}
	}

	var out strings.Builder

	fmt.Fprintf(&out, "<memory-digest>\n%s\nProject key: %s (save project notes in `%s`).\n", Guidance, key, ProjectDir(key))

	section(&out, "Project notes", project)
	section(&out, "User and feedback notes", global)

	out.WriteString("</memory-digest>")

	return out.String(), nil
}

// listing is a digest line and the date its note was created.
type listing struct {
	date string
	line string
}

// section lists notes under a heading, newest first, stopping before the digest outgrows its budget.
func section(out *strings.Builder, heading string, lines []listing) {
	if len(lines) == 0 {
		return
	}

	slices.SortFunc(lines, func(a, b listing) int {
		return cmp.Or(strings.Compare(b.date, a.date), strings.Compare(a.line, b.line))
	})

	fmt.Fprintf(out, "\n%s:\n", heading)

	var size int

	for i, l := range lines {
		size += len(l.line) + 1
		if size > digestSectionMax {
			fmt.Fprintf(out, "- ... %d more, use semantic_search_notes\n", len(lines)-i)

			return
		}

		out.WriteString(l.line + "\n")
	}
}

// inProject returns whether the path is a note in the project's directory.
func inProject(path, key string) bool {
	return strings.HasPrefix(path, ProjectDir(key)+"/")
}
