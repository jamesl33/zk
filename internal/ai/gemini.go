package ai

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/jamesl33/zk/internal/ai/cache"
	"google.golang.org/genai"
)

// Gemini defines a client for interacting with the Gemini API.
//
// TODO (jamesl33): Make the model configurable.
// TODO (jamesl33): Handle the 8k context window.
type Gemini struct {
	client *genai.Client
	gcache *cache.Cache[string]
	ecache *cache.Cache[[]byte]
}

var _ Client = (*Gemini)(nil)

// New creates a new client for interacting with the Gemini API.
func New(ctx context.Context, path string) (*Gemini, error) {
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
		client: ai,
		gcache: gca,
		ecache: eca,
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

	resp, err := g.client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, &genai.GenerateContentConfig{})
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

	resp, err := g.client.Models.EmbedContent(ctx, "gemini-embedding-2-preview", contents, &genai.EmbedContentConfig{})
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
