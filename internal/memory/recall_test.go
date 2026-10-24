package memory

import (
	"context"
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/jamesl33/zk/internal/note"
	"github.com/stretchr/testify/require"
)

const prompt = "how do I run the tests here"

func TestRecallInjectsOnceThenDedups(t *testing.T) {
	vault := withVault(t)
	path := writeNote(t, "3 Resources/a.md", "Tests need -shuffle", []string{"reference"}, "Run go test with -shuffle=on.")

	find := finderOf(t, 0.5, path)

	got, err := Recall(context.Background(), vault, "s", prompt, find, MaxDistance)
	require.NoError(t, err)
	require.Contains(t, got, "<memory>")
	require.Contains(t, got, "Tests need -shuffle (3 Resources/a.md, saved 2026-01-02)")
	require.Contains(t, got, "Run go test with -shuffle=on.")

	got, err = Recall(context.Background(), vault, "s", prompt, find, MaxDistance)
	require.NoError(t, err)
	require.Empty(t, got)

	// Compaction drops injected notes from context, so they may be injected again.
	require.NoError(t, ResetRecalled(vault, "s"))

	got, err = Recall(context.Background(), vault, "s", prompt, find, MaxDistance)
	require.NoError(t, err)
	require.Contains(t, got, "Tests need -shuffle")
}

func TestRecallThresholdIsPassedToTheFinder(t *testing.T) {
	vault := withVault(t)
	path := writeNote(t, "3 Resources/a.md", "A", nil, "body")

	got, err := Recall(context.Background(), vault, "s", prompt, finderOf(t, 0.5, path), 0.4)
	require.NoError(t, err)
	require.Empty(t, got)
}

func TestRecallSkips(t *testing.T) {
	vault := withVault(t)
	path := writeNote(t, "3 Resources/a.md", "A", nil, "body")

	for _, p := range []string{"", "yes", "continue please", "/clear all of the things"} {
		got, err := Recall(context.Background(), vault, "s", p, finderOf(t, 0.1, path), MaxDistance)
		require.NoError(t, err)
		require.Empty(t, got, p)
	}
}

func TestRecallIncludesOtherProjects(t *testing.T) {
	vault := withVault(t)
	mine := writeNote(t, "1 Projects/mine/a.md", "Mine", nil, "mine")
	other := writeNote(t, "1 Projects/other/b.md", "Other", nil, "other")

	got, err := Recall(context.Background(), vault, "s", prompt, finderOf(t, 0.1, mine, other), MaxDistance)
	require.NoError(t, err)
	require.Contains(t, got, "Mine")
	require.Contains(t, got, "Other")
}

func TestRecallCapsSize(t *testing.T) {
	vault := withVault(t)
	big := string(make([]byte, recallMax))
	a := writeNote(t, "3 Resources/a.md", "A", nil, "small")
	b := writeNote(t, "3 Resources/b.md", "B", nil, "x"+big)

	got, err := Recall(context.Background(), vault, "s", prompt, finderOf(t, 0.1, a, b), MaxDistance)
	require.NoError(t, err)
	require.Contains(t, got, "small")
	require.Less(t, len(got), recallMax+500)
}

func TestRecallFinderError(t *testing.T) {
	vault := withVault(t)

	fail := func(context.Context, string, int, float32) ([]*note.Note, error) {
		return nil, errors.New("ollama down")
	}

	_, err := Recall(context.Background(), vault, "s", prompt, fail, MaxDistance)
	require.ErrorContains(t, err, "ollama down")
}

func TestRecallTruncatesOversizedNote(t *testing.T) {
	vault := withVault(t)
	big := writeNote(t, "3 Resources/big.md", "Big", nil, strings.Repeat("é", recallMax))

	got, err := Recall(context.Background(), vault, "s", prompt, finderOf(t, 0.1, big), MaxDistance)
	require.NoError(t, err)
	require.Contains(t, got, "### Big")
	require.Contains(t, got, `full note: read_note "3 Resources/big.md"`)
	require.True(t, utf8.ValidString(got))
	require.Less(t, len(got), recallMax+500)
}
