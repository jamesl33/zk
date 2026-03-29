package ai

import "context"

//go:generate mockgen -destination mocks/mocks.go github.com/jamesl33/zk/internal/ai Client

// Client provides an API for performing AI operations.
type Client interface {
	// Generate text from a prompt.
	Generate(ctx context.Context, prompt string) (string, error)

	// Embed a string of text.
	Embed(ctx context.Context, content string) ([]float32, error)
}
