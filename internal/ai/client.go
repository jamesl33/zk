package ai

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/jamesl33/zk/internal/ai/cache"
	"google.golang.org/genai"
)

// Client - TODO
//
// TODO (jamesl33): Make the model configurable.
// TODO (jamesl33): Handle the 2k context window.
type Client struct {
	client *genai.Client
	gcache *cache.Cache[string]
	ecache *cache.Cache[[]byte]
}

// New - TODO
func New(ctx context.Context, path string) (*Client, error) {
	ai, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	gca, err := cache.New[string](ctx, path, "cache_generate")
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	eca, err := cache.New[[]byte](ctx, path, "cache_embed")
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	client := Client{
		client: ai,
		gcache: gca,
		ecache: eca,
	}

	return &client, nil
}

// Generate - TODO
func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	cached, err := c.gcache.Get(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("%w", err) // TODO
	}

	// TODO
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
		return "", fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if len(resp.Candidates) != 1 || len(resp.Candidates[0].Content.Parts) != 1 {
		return "", nil
	}

	// TODO
	result := resp.Candidates[0].Content.Parts[0].Text

	// TODO
	err = c.gcache.Set(ctx, prompt, result)
	if err != nil {
		return "", fmt.Errorf("%w", err) // TODO
	}

	return result, nil
}

// Embed - TODO
func (c *Client) Embed(ctx context.Context, content string) ([]float32, error) {
	cached, err := c.ecache.Get(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// TODO
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
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	if len(resp.Embeddings) != 1 {
		return nil, nil
	}

	// TODO
	result := resp.Embeddings[0].Values

	// TODO
	blob, err := sf32toblob(result)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	// TODO
	err = c.ecache.Set(ctx, content, blob)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return result, nil
}

// sf32toblob - TODO
func sf32toblob(embedding []float32) ([]byte, error) {
	blob := make([]byte, 4*len(embedding))

	_, err := binary.Encode(blob, binary.LittleEndian, embedding)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return blob, nil
}

// blobtosf32 - TODO
func blobtosf32(blob []byte) ([]float32, error) {
	embedding := make([]float32, len(blob)/4)

	_, err := binary.Decode(blob, binary.LittleEndian, embedding)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	return embedding, nil
}
