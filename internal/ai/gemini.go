package ai

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/jamesl33/zk/internal/ai/cache"
	"google.golang.org/genai"
)

// maxContextTokens is the approximate context window supported by the embedding model; inputs
// larger than this are rejected up front rather than sent to the API to fail. Notes are expected
// to stay small/atomic (Zettelkasten style), so this limit should rarely be hit in practice.
const maxContextTokens = 8000

// Gemini defines a client for interacting with the Gemini API.
type Gemini struct {
	client    *genai.Client
	gcache    *cache.Cache[string]
	ecache    *cache.Cache[[]byte]
	model     string
	embedding string
}

var _ Client = (*Gemini)(nil)

// Option configures a Gemini client.
type Option func(*Gemini)

// WithModel overrides the model used for text generation.
func WithModel(model string) Option {
	return func(g *Gemini) { g.model = model }
}

// WithEmbedModel overrides the model used for embedding.
func WithEmbedModel(model string) Option {
	return func(g *Gemini) { g.embedding = model }
}

// NewGemini creates a new client for interacting with the Gemini API.
func NewGemini(ctx context.Context, path string, opts ...Option) (*Gemini, error) {
	ai, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	gca, err := cache.New[string](ctx, path, "cache_generate")
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	eca, err := cache.New[[]byte](ctx, path, "cache_embed")
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

	client := Gemini{
		client:    ai,
		gcache:    gca,
		ecache:    eca,
		model:     "gemini-2.5-flash",
		embedding: "gemini-embedding-2-preview",
	}

	// The model is deliberately taken from the environment rather than a per-call argument: the
	// response/embedding cache keys on the prompt/content alone, so switching models between
	// invocations against the same vault would silently serve stale results from a different
	// model. An Option can still override it explicitly (e.g. for tests).
	if model := os.Getenv("ZK_GEMINI_MODEL"); model != "" {
		client.model = model
	}

	if model := os.Getenv("ZK_GEMINI_EMBED_MODEL"); model != "" {
		client.embedding = model
	}

	for _, opt := range opts {
		opt(&client)
	}

	return &client, nil
}

// Generate text from a prompt.
func (g *Gemini) Generate(ctx context.Context, prompt string) (string, error) {
	cached, err := g.gcache.Get(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to get from cache: %w", err)
	}

	if cached != nil {
		return *cached, nil
	}

	part := &genai.Part{
		Text: prompt,
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{part}},
	}

	if err := g.checkSize(ctx, g.model, contents); err != nil {
		return "", err
	}

	resp, err := g.client.Models.GenerateContent(ctx, g.model, contents, &genai.GenerateContentConfig{})
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) != 1 || len(resp.Candidates[0].Content.Parts) != 1 {
		return "", nil
	}

	result := resp.Candidates[0].Content.Parts[0].Text

	err = g.gcache.Set(ctx, prompt, result)
	if err != nil {
		return "", fmt.Errorf("failed to set cache: %w", err)
	}

	return result, nil
}

// Embed a string of text.
func (g *Gemini) Embed(ctx context.Context, content string) ([]float32, error) {
	cached, err := g.ecache.Get(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	if cached != nil {
		return blobtosf32(*cached)
	}

	part := &genai.Part{
		Text: content,
	}

	contents := []*genai.Content{
		{Parts: []*genai.Part{part}},
	}

	if err := g.checkSize(ctx, g.embedding, contents); err != nil {
		return nil, err
	}

	resp, err := g.client.Models.EmbedContent(ctx, g.embedding, contents, &genai.EmbedContentConfig{})
	if err != nil {
		return nil, fmt.Errorf("failed to embed content: %w", err)
	}

	if len(resp.Embeddings) != 1 {
		return nil, nil
	}

	result := resp.Embeddings[0].Values

	blob, err := sf32toblob(result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to blob: %w", err)
	}

	err = g.ecache.Set(ctx, content, blob)
	if err != nil {
		return nil, fmt.Errorf("failed to set cache: %w", err)
	}

	return result, nil
}

// checkSize returns an error if the given contents exceed the supported context window.
func (g *Gemini) checkSize(ctx context.Context, model string, contents []*genai.Content) error {
	resp, err := g.client.Models.CountTokens(ctx, model, contents, &genai.CountTokensConfig{})
	if err != nil {
		return fmt.Errorf("failed to count tokens: %w", err)
	}

	if resp.TotalTokens > maxContextTokens {
		return fmt.Errorf("input is too large (%d tokens, max %d)", resp.TotalTokens, maxContextTokens)
	}

	return nil
}

// sf32toblob converts a slice of float32 to a blob.
func sf32toblob(embedding []float32) ([]byte, error) {
	blob := make([]byte, 4*len(embedding))

	_, err := binary.Encode(blob, binary.LittleEndian, embedding)
	if err != nil {
		return nil, fmt.Errorf("failed to encode embedding: %w", err)
	}

	return blob, nil
}

// blobtosf32 converts a blob to a slice of float32.
func blobtosf32(blob []byte) ([]float32, error) {
	embedding := make([]float32, len(blob)/4)

	_, err := binary.Decode(blob, binary.LittleEndian, embedding)
	if err != nil {
		return nil, fmt.Errorf("failed to decode embedding: %w", err)
	}

	return embedding, nil
}
