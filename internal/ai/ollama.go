package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/jamesl33/zk/internal/ai/cache"
	"github.com/jamesl33/zk/internal/ptr"
	"github.com/ollama/ollama/api"
)

// Ollama defines a client for embedding text against a local Ollama instance. Embedding is called
// far more often than generation (every note, plus every vector search query), so it's kept local
// to stay fast/free rather than paying hosted API cost per call.
type Ollama struct {
	client    *api.Client
	ecache    *cache.Cache[[]byte]
	embedding string
}

var _ Client = (*Ollama)(nil)

// OllamaOption configures an Ollama client.
type OllamaOption func(*Ollama)

// WithOllamaEmbedModel overrides the model used for embedding.
func WithOllamaEmbedModel(model string) OllamaOption {
	return func(o *Ollama) { o.embedding = model }
}

// NewOllama creates a new client for embedding text against a local Ollama instance. The instance
// is located via Ollama's own OLLAMA_HOST environment variable (see api.ClientFromEnvironment).
func NewOllama(ctx context.Context, path string, opts ...OllamaOption) (*Ollama, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	eca, err := cache.New[[]byte](ctx, path, "cache_embed")
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	ollama := Ollama{
		client:    client,
		ecache:    eca,
		embedding: "embeddinggemma",
	}

	if model := os.Getenv("ZK_OLLAMA_EMBED_MODEL"); model != "" {
		ollama.embedding = model
	}

	for _, opt := range opts {
		opt(&ollama)
	}

	return &ollama, nil
}

// Generate is unimplemented; Ollama is only ever used for Embed.
func (o *Ollama) Generate(_ context.Context, _ string) (string, error) {
	panic("unimplemented")
}

// Embed a string of text.
func (o *Ollama) Embed(ctx context.Context, content string) ([]float32, error) {
	cached, err := o.ecache.Get(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	if cached != nil {
		return blobtosf32(*cached)
	}

	resp, err := o.client.Embed(ctx, &api.EmbedRequest{
		Model:    o.embedding,
		Input:    content,
		Truncate: ptr.To(false),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to embed content: %w", err)
	}

	if len(resp.Embeddings) != 1 {
		return nil, nil
	}

	result := resp.Embeddings[0]

	blob, err := sf32toblob(result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to blob: %w", err)
	}

	err = o.ecache.Set(ctx, content, blob)
	if err != nil {
		return nil, fmt.Errorf("failed to set cache: %w", err)
	}

	return result, nil
}
