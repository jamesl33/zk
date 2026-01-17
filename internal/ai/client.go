package ai

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/jamesl33/zk/internal/ai/cache"
	"google.golang.org/genai"
)

// Client defines a client for interacting with the Gemini API.
//
// TODO (jamesl33): Make the model configurable.
// TODO (jamesl33): Handle the 2k context window.
type Client struct {
	client *genai.Client
	gcache *cache.Cache[string]
	ecache *cache.Cache[[]byte]
}

// New creates a new client for interacting with the Gemini API.
func New(ctx context.Context, path string) (*Client, error) {
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

	client := Client{
		client: ai,
		gcache: gca,
		ecache: eca,
	}

	return &client, nil
}

// Generate generates text from a prompt.
func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	cached, err := c.gcache.Get(ctx, prompt)
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

	resp, err := c.client.Models.GenerateContent(ctx, "gemini-2.5-flash", contents, &genai.GenerateContentConfig{})
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) != 1 || len(resp.Candidates[0].Content.Parts) != 1 {
		return "", nil
	}

	result := resp.Candidates[0].Content.Parts[0].Text

	err = c.gcache.Set(ctx, prompt, result)
	if err != nil {
		return "", fmt.Errorf("failed to set cache: %w", err)
	}

	return result, nil
}

// Embed creates an embedding from a string of text.
func (c *Client) Embed(ctx context.Context, content string) ([]float32, error) {
	cached, err := c.ecache.Get(ctx, content)
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

	resp, err := c.client.Models.EmbedContent(ctx, "gemini-embedding-001", contents, &genai.EmbedContentConfig{})
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

	err = c.ecache.Set(ctx, content, blob)
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
