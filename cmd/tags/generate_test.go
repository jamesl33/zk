package tags

import (
	"os"
	"path/filepath"
	"testing"

	mock_ai "github.com/jamesl33/zk/internal/ai/mocks"
	"github.com/jamesl33/zk/internal/note"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGenerateGenerate(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\n---\nBody about hiking"), 0o644)
	require.NoError(t, err)

	n, err := note.New(path)
	require.NoError(t, err)

	mclient.
		EXPECT().
		Generate(gomock.Any(), gomock.Any()).
		Return("```yaml\ntags:\n  - thru hiking\n```", nil)

	g := Generate{}

	err = g.generate(t.Context(), mclient, n)
	require.NoError(t, err)
	assert.Equal(t, []string{"thru_hiking"}, n.Frontmatter.Tags)
}

func TestGenerateGenerateEmptyBody(t *testing.T) {
	var (
		tmp     = t.TempDir()
		ctrl    = gomock.NewController(t)
		mclient = mock_ai.NewMockClient(ctrl)
	)

	path := filepath.Join(tmp, "note.md")

	err := os.WriteFile(path, []byte("---\ntitle: Note 1\n---\n"), 0o644)
	require.NoError(t, err)

	n, err := note.New(path)
	require.NoError(t, err)

	// No calls to the client are expected for an empty body.
	g := Generate{}

	err = g.generate(t.Context(), mclient, n)
	require.NoError(t, err)
}
