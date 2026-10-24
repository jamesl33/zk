package memory

import (
	"context"
	"fmt"
	"strings"

	"github.com/jamesl33/zk/internal/note"
)

const (
	// recallLimit is the most notes searched for per prompt.
	recallLimit = 15

	// recallMax is the most bytes of notes injected per prompt.
	recallMax = 2000

	// MaxDistance is the default cosine distance beyond which a note is not relevant, so that trivial prompts inject
	// nothing. It's overridden by ZK_MEMORY_MAX_DISTANCE.
	MaxDistance float32 = 0.65

	// minWords is the shortest prompt worth searching for ("yes", "continue" can't match anything).
	minWords = 3

	// queryMax is the most bytes of the prompt which is embedded.
	queryMax = 2000
)

// Recall finds the notes relevant to a prompt, which haven't already been injected into the session. It runs in the vault
// and returns nothing when there's nothing worth saying.
func Recall(ctx context.Context, vault, session, prompt string, find Finder, maxDistance float32) (string, error) {
	if strings.HasPrefix(prompt, "/") || len(strings.Fields(prompt)) < minWords {
		return "", nil
	}

	notes, err := find(ctx, truncate(prompt, queryMax), recallLimit, maxDistance)
	if err != nil {
		return "", fmt.Errorf("failed to find notes: %w", err)
	}

	recalled, err := Recalled(vault, session)
	if err != nil {
		return "", fmt.Errorf("failed to read recalled notes: %w", err)
	}

	var (
		out   strings.Builder
		paths []string
	)

	for _, n := range notes {
		if recalled[n.Path] {
			continue
		}

		entry, err := format(n)
		if err != nil {
			return "", err
		}

		if out.Len()+len(entry) > recallMax {
			break
		}

		out.WriteString(entry)

		paths = append(paths, n.Path)
	}

	if len(paths) == 0 {
		return "", nil
	}

	err = MarkRecalled(vault, session, paths...)
	if err != nil {
		return "", fmt.Errorf("failed to record recalled notes: %w", err)
	}

	return "<memory>\nNotes from Claude's memory vault that may be relevant. Older notes may be stale: verify anything about code or state before relying on it. If a note is wrong or outdated, fix or delete it with the memory tools.\n\n" + out.String() + "</memory>", nil
}

// format renders a note for injection.
func format(n *note.Note) (string, error) {
	body, err := n.GetBody()
	if err != nil {
		return "", fmt.Errorf("failed to read note: %w", err)
	}

	return fmt.Sprintf("### %s (%s, saved %s)\n%s\n\n", n.Frontmatter.Title, n.Path, n.Frontmatter.Date, strings.TrimSpace(body)), nil
}

// truncate shortens s to at most n bytes.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[:n]
}
