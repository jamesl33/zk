package memory

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func userLine(text string) string {
	return `{"type":"user","message":{"content":"` + text + `"}}`
}

func writeTranscript(t *testing.T, lines ...string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "transcript.jsonl")

	require.NoError(t, os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644))

	return path
}

// recorder is a Runner which records its prompts, failing the nth call if fail is set.
type recorder struct {
	prompts []string
	fail    int
	effect  func(vault string, call int)
}

func (r *recorder) run(_ context.Context, vault, prompt string) error {
	r.prompts = append(r.prompts, prompt)

	call := len(r.prompts)
	if r.fail == call {
		return errors.New("boom")
	}

	if r.effect != nil {
		r.effect(vault, call)
	}

	return nil
}

func TestExtractSavesAdvancesOffsetAndCommits(t *testing.T) {
	vault := withVault(t)

	transcript := writeTranscript(t, userLine("remember I like tabs"), `{"type":"summary"}`, userLine("ok"))

	r := recorder{effect: func(vault string, _ int) {
		require.NoError(t, os.WriteFile(filepath.Join(vault, "3 Resources", "n.md"), []byte("x"), 0o644))
	}}

	err := Extract(context.Background(), Extraction{Vault: vault, Session: "s", Transcript: transcript, Key: "p", Run: r.run})
	require.NoError(t, err)

	require.Len(t, r.prompts, 1)
	require.Contains(t, r.prompts[0], "user: remember I like tabs")
	require.Contains(t, r.prompts[0], "project \"p\"")
	require.NotContains(t, r.prompts[0], "summary")
	require.Equal(t, 3, Offset(vault, "s"))

	out, err := exec.Command("git", "-C", vault, "log", "--format=%s").Output()
	require.NoError(t, err)
	require.Equal(t, "Save memories\n", string(out))
}

func TestExtractOnlySeesNewLines(t *testing.T) {
	vault := withVault(t)

	transcript := writeTranscript(t, userLine("alpha"), userLine("beta"))

	require.NoError(t, SetOffset(vault, "s", 1))

	r := recorder{}

	require.NoError(t, Extract(context.Background(), Extraction{Vault: vault, Session: "s", Transcript: transcript, Run: r.run}))
	require.Len(t, r.prompts, 1)
	require.NotContains(t, r.prompts[0], "alpha")
	require.Contains(t, r.prompts[0], "beta")
}

func TestExtractFailureKeepsOffset(t *testing.T) {
	vault := withVault(t)

	r := recorder{fail: 1}

	err := Extract(context.Background(), Extraction{Vault: vault, Session: "s", Transcript: writeTranscript(t, userLine("x")), Run: r.run})
	require.ErrorContains(t, err, "boom")
	require.Zero(t, Offset(vault, "s"))
}

func TestExtractOffsetAdvancesPerChunk(t *testing.T) {
	vault := withVault(t)

	big := strings.Repeat("a", chunkSize*6/10)
	transcript := writeTranscript(t, userLine(big), userLine(big), userLine(big))

	r := recorder{fail: 2}

	err := Extract(context.Background(), Extraction{Vault: vault, Session: "s", Transcript: transcript, Run: r.run})
	require.Error(t, err)
	require.Equal(t, 1, Offset(vault, "s")) // The first chunk was saved, the retry starts at the second
}

func TestExtractDropsOldestChunksPastTheCap(t *testing.T) {
	vault := withVault(t)

	big := strings.Repeat("a", chunkSize*6/10)

	var lines []string
	for range maxChunks + 2 {
		lines = append(lines, userLine(big))
	}

	r := recorder{}

	require.NoError(t, Extract(context.Background(), Extraction{Vault: vault, Session: "s", Transcript: writeTranscript(t, lines...), Run: r.run}))
	require.Len(t, r.prompts, maxChunks)
	require.Equal(t, maxChunks+2, Offset(vault, "s"))
}

func TestExtractFixesLintWithTheRealReport(t *testing.T) {
	vault := withVault(t)

	r := recorder{effect: func(vault string, call int) {
		if call == 1 { // An orphan: nothing links to or from it
			writeNote(t, filepath.Join(vault, "3 Resources", "20260102030405.md"), "Orphan", nil, "no links")
		}
	}}

	err := Extract(context.Background(), Extraction{Vault: vault, Session: "s", Transcript: writeTranscript(t, userLine("x")), Run: r.run})
	require.NoError(t, err)

	require.Len(t, r.prompts, 2)
	require.Contains(t, r.prompts[1], "20260102030405.md")
}

func TestRenderKeepsToolCallsAndResultsClipped(t *testing.T) {
	long := strings.Repeat("z", clipSize*2)

	got := render([]byte(`{"type":"assistant","message":{"content":[{"type":"text","text":"hi"},{"type":"tool_use","name":"Bash","input":{"command":"ls"}},{"type":"tool_result","content":[{"type":"text","text":"` + long + `"}]}]}}`))

	require.Contains(t, got, "assistant: hi")
	require.Contains(t, got, `[tool Bash] {"command":"ls"}`)
	require.Contains(t, got, "[result] "+strings.Repeat("z", clipSize)+"…")
	require.NotContains(t, got, strings.Repeat("z", clipSize+1))
}

func TestSavingSuppressedByEnv(t *testing.T) {
	t.Setenv(SavingEnv, "")
	require.False(t, Saving())

	t.Setenv(SavingEnv, "1")
	require.True(t, Saving())
}

func TestEnsureVaultIsIdempotent(t *testing.T) {
	vault := filepath.Join(t.TempDir(), "v")

	require.NoError(t, EnsureVault(vault))
	require.NoError(t, EnsureVault(vault))

	for _, f := range []string{".git", ".zk/state", "ZK.md", "CLAUDE.md", "1 Projects", ".gitignore"} {
		_, err := os.Stat(filepath.Join(vault, f))
		require.NoError(t, err, f)
	}
}
