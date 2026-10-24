package memory

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/jamesl33/zk/internal/lister"
)

// digestMax is the most bytes of note titles in the digest.
const digestMax = 2000

// globalTags are the tags of notes which are always listed, whichever project they were saved from.
var globalTags = []string{"user", "feedback"}

// Digest is what Claude is told at the start of a session (and again after compaction): how to save memories, and the
// titles of those worth knowing before the first prompt. It runs in the vault.
func Digest(ctx context.Context, key string) (string, error) {
	lister, err := lister.NewLister(lister.WithPath("."))
	if err != nil {
		return "", fmt.Errorf("failed to create lister: %w", err)
	}

	var project, global []string

	for n, err := range lister.Many(ctx) {
		if err != nil {
			return "", fmt.Errorf("failed to list notes: %w", err)
		}

		line := fmt.Sprintf("- %s (%s)", n.Frontmatter.Title, n.Path)

		switch {
		case inProject(n.Path, key):
			project = append(project, line)
		case slices.ContainsFunc(n.Frontmatter.Tags, func(tag string) bool { return slices.Contains(globalTags, tag) }):
			global = append(global, line)
		}
	}

	var out strings.Builder

	fmt.Fprintf(&out, "<memory-digest>\n%s\nProject key: %s (save project notes in `%s`).\n", Guidance, key, ProjectDir(key))

	section(&out, "Project notes", project)
	section(&out, "User and feedback notes", global)

	out.WriteString("</memory-digest>")

	return out.String(), nil
}

// section lists lines under a heading, stopping before the digest outgrows its budget.
func section(out *strings.Builder, heading string, lines []string) {
	if len(lines) == 0 {
		return
	}

	slices.Sort(lines)

	fmt.Fprintf(out, "\n%s:\n", heading)

	var size int

	for i, line := range lines {
		size += len(line) + 1
		if size > digestMax/2 {
			fmt.Fprintf(out, "- ... %d more, use semantic_search_notes\n", len(lines)-i)

			return
		}

		out.WriteString(line + "\n")
	}
}

// inProject returns whether the path is a note in the project's directory.
func inProject(path, key string) bool {
	return strings.HasPrefix(path, ProjectDir(key)+"/")
}
