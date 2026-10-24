package memory

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDigest(t *testing.T) {
	withVault(t)

	writeNote(t, "1 Projects/mine/a.md", "Mine note", []string{"project"}, "x")
	writeNote(t, "1 Projects/other/b.md", "Other project note", []string{"project"}, "x")
	writeNote(t, "3 Resources/c.md", "Prefers tabs", []string{"user"}, "x")
	writeNote(t, "3 Resources/d.md", "Plain fact", []string{"reference"}, "x")

	got, err := Digest(context.Background(), "mine")
	require.NoError(t, err)

	require.Contains(t, got, Guidance)
	require.Contains(t, got, "Project key: mine (save project notes in `1 Projects/mine`)")
	require.Contains(t, got, "- Mine note (1 Projects/mine/a.md)")
	require.Contains(t, got, "- Prefers tabs (3 Resources/c.md)")
	require.NotContains(t, got, "Other project note")
	require.NotContains(t, got, "Plain fact")
}

func TestDigestCapsTitles(t *testing.T) {
	withVault(t)

	for i := range 100 {
		writeNote(t, fmt.Sprintf("1 Projects/mine/%03d.md", i), fmt.Sprintf("A fairly long note title number %03d", i), nil, "x")
	}

	got, err := Digest(context.Background(), "mine")
	require.NoError(t, err)
	require.Contains(t, got, "more, use semantic_search_notes")
	require.Less(t, len(got), len(Guidance)+digestMax+300)
}
