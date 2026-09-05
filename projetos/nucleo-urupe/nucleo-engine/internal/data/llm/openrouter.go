/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenRouterProvider struct {
	config Config
	http   *http.Client
}

func NewOpenRouterProvider(config Config) *OpenRouterProvider {
	if config.OpenRouterBaseURL == "" {
		config.OpenRouterBaseURL = "https://openrouter.ai/api/v1"
	}
	return &OpenRouterProvider{
		config: config,
		http: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *OpenRouterProvider) Complete(ctx context.Context, messages []Message, stop []string, opts RequestOptions) (Completion, error) {
	reqBody := ChatRequest{
		Model:       opts.Model,
		Messages:    messages,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
		Stop:        stop,
		Tools:       opts.Tools,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return Completion{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(p.config.OpenRouterBaseURL, "/"))
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(data))
	if err != nil {
		return Completion{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.OpenRouterAPIKey)
	if p.config.OpenRouterSiteURL != "" {
		req.Header.Set("HTTP-Referer", p.config.OpenRouterSiteURL)
	}
	if p.config.OpenRouterTitle != "" {
		req.Header.Set("X-OpenRouter-Title", p.config.OpenRouterTitle)
	}

	resp, err := p.http.Do(req)
	if err != nil {
		return Completion{}, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Completion{}, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return Completion{}, fmt.Errorf("rate limited (429)")
	}

	if resp.StatusCode != http.StatusOK {
		return Completion{}, fmt.Errorf("openrouter api error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return Completion{}, fmt.Errorf("failed to decode response: %w", err)
	}

	if chatResp.Error != nil {
		return Completion{}, fmt.Errorf("openrouter api error: %s", chatResp.Error.Message)
	}
	if len(chatResp.Choices) == 0 {
		return Completion{}, fmt.Errorf("no content returned from openrouter")
	}

	return Completion{
		Text:      chatResp.Choices[0].Message.Content,
		Reasoning: chatResp.Choices[0].Message.Reasoning,
		Model:     chatResp.Model,
		Usage:     chatResp.Usage,
		ToolCalls: chatResp.Choices[0].Message.ToolCalls,
	}, nil
}

func (p *OpenRouterProvider) Embed(ctx context.Context, input string) ([]float32, error) {
	reqBody := struct {
		Model string   `json:"model"`
		Input []string `json:"input"`
	}{
		Model: "text-embedding-3-small",
		Input: []string{input},
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/embeddings", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.OpenRouterAPIKey)

	resp, err := p.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var embResp struct {
		Data []struct {
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyBytes, &embResp); err != nil {
		return nil, err
	}

	if len(embResp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return embResp.Data[0].Embedding, nil
}
