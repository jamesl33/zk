package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/jamesl33/zk/internal/ai/cache"
	"google.golang.org/genai"
)

// Gemini defines a client for interacting with the Gemini API.
type Gemini struct {
	client *genai.Client
	gcache *cache.Cache[string]
	model  string
}

var _ Client = (*Gemini)(nil)

// Option configures a Gemini client.
type Option func(*Gemini)

// WithModel overrides the model used for text generation.
func WithModel(model string) Option {
	return func(g *Gemini) { g.model = model }
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

	client := Gemini{
		client: ai,
		gcache: gca,
		model:  "gemini-2.5-flash",
	}

	// The model is deliberately taken from the environment rather than a per-call argument: the response cache keys on
	// the prompt alone, so switching models between invocations against the same vault would silently serve stale
	// results from a different model. An Option can still override it explicitly (e.g. for tests).
	if model := os.Getenv("ZK_GEMINI_MODEL"); model != "" {
		client.model = model
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

// Embed is unimplemented; Gemini is only ever used for Generate.
func (g *Gemini) Embed(_ context.Context, _ string) ([]float32, error) {
	panic("unimplemented")
}
