package ai

import "context"

//go:generate go run go.uber.org/mock/mockgen -destination mocks/mocks.go github.com/jamesl33/zk/internal/ai Embedder,Generator

// Embedder provides an API for embedding text.
type Embedder interface {
	// Embed a string of text.
	Embed(ctx context.Context, content string) ([]float32, error)
}

// Generator provides an API for generating text.
type Generator interface {
	// Generate text from a prompt.
	Generate(ctx context.Context, prompt string) (string, error)
}
