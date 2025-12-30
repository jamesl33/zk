package ai

import (
	"context"
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
	cache  *cache.Cache
}

// New - TODO
func New(ctx context.Context, path string) (*Client, error) {
	ai, err := genai.NewClient(ctx, &genai.ClientConfig{
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	ca, err := cache.New(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("%w", err) // TODO
	}

	client := Client{
		client: ai,
		cache:  ca,
	}

	return &client, nil
}

// Generate - TODO
func (c *Client) Generate(ctx context.Context, prompt string) (string, error) {
	cached, err := c.cache.Get(ctx, prompt)
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
	err = c.cache.Set(ctx, prompt, result)
	if err != nil {
		return "", fmt.Errorf("%w", err) // TODO
	}

	return result, nil
}

// Embed - TODO
func (c *Client) Embed(ctx context.Context, content string) ([]float32, error) {
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

	return resp.Embeddings[0].Values, nil
}
