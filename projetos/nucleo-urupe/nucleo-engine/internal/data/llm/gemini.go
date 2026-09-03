/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package llm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"google.golang.org/genai"
)

type GeminiProvider struct {
	client *genai.Client
}

func NewGeminiProvider(ctx context.Context, apiKey string) (*GeminiProvider, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}
	return &GeminiProvider{client: client}, nil
}

func (g *GeminiProvider) Complete(ctx context.Context, messages []Message, stop []string, opts RequestOptions) (Completion, error) {
	modelName := opts.Model
	if !strings.HasPrefix(modelName, "models/") && !strings.Contains(modelName, "/") {
		modelName = "models/" + modelName
	}

	config := &genai.GenerateContentConfig{
		MaxOutputTokens: int32(opts.MaxTokens),
	}
	if opts.HasTemperature {
		temp := float32(opts.Temperature)
		config.Temperature = &temp
	}
	if len(stop) > 0 {
		config.StopSequences = stop
	}

	if opts.ResponseMimeType != "" {
		config.ResponseMIMEType = opts.ResponseMimeType
	}
	if opts.ResponseSchema != nil {
		if schema, ok := opts.ResponseSchema.(*genai.Schema); ok {
			config.ResponseSchema = schema
		}
	}

	if opts.ServiceTier != "" {
		config.ServiceTier = genai.ServiceTier(opts.ServiceTier)
	}

	// Apply robust safety settings for an unfiltered agent persona
	config.SafetySettings = []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockThresholdBlockNone},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockThresholdBlockNone},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockThresholdBlockNone},
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockThresholdBlockNone},
	}

	if len(messages) == 0 {
		return Completion{}, fmt.Errorf("no messages provided")
	}

	var contents []*genai.Content
	for _, m := range messages {
		if m.Role == "system" {
			sysContent, err := g.convertToGeminiContent(ctx, m)
			if err == nil {
				config.SystemInstruction = sysContent
			}
			continue
		}
		content, err := g.convertToGeminiContent(ctx, m)
		if err != nil {
			return Completion{}, err
		}
		contents = append(contents, content)
	}

	resp, err := g.client.Models.GenerateContent(ctx, modelName, contents, config)
	if err != nil {
		return Completion{}, err
	}

	if len(resp.Candidates) == 0 {
		return Completion{}, fmt.Errorf("no candidates returned from gemini")
	}

	usage := Usage{}
	if resp.UsageMetadata != nil {
		usage.PromptTokens = int(resp.UsageMetadata.PromptTokenCount)
		usage.CompletionTokens = int(resp.UsageMetadata.CandidatesTokenCount)
		usage.TotalTokens = int(resp.UsageMetadata.TotalTokenCount)
	}

	return Completion{
		Text:  resp.Text(),
		Model: modelName,
		Usage: usage,
	}, nil
}

func (g *GeminiProvider) convertToGeminiContent(ctx context.Context, msg Message) (*genai.Content, error) {
	role := "user"
	if msg.Role == "assistant" || msg.Role == "bot" || msg.Role == "model" {
		role = "model"
	}

	content := &genai.Content{Role: role}

	switch c := msg.Content.(type) {
	case string:
		content.Parts = append(content.Parts, &genai.Part{Text: c})
	case []ContentPart:
		for _, part := range c {
			if part.Type == "text" {
				content.Parts = append(content.Parts, &genai.Part{Text: part.Text})
			} else if part.Type == "image" && len(part.Data) > 0 {
				content.Parts = append(content.Parts, &genai.Part{
					InlineData: &genai.Blob{
						MIMEType: part.MimeType,
						Data:     part.Data,
					},
				})
			} else if part.Type == "image_url" && part.ImageURL != nil {
				data, mime, err := g.downloadImage(ctx, part.ImageURL.URL)
				if err != nil {
					return nil, fmt.Errorf("failed to download image %s: %w", part.ImageURL.URL, err)
				}
				content.Parts = append(content.Parts, &genai.Part{
					InlineData: &genai.Blob{
						MIMEType: mime,
						Data:     data,
					},
				})
			}
		}
	default:
		return nil, fmt.Errorf("unsupported content type: %T", msg.Content)
	}

	return content, nil
}

func (g *GeminiProvider) downloadImage(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to fetch image: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	mime := resp.Header.Get("Content-Type")
	if mime == "" {
		mime = "image/jpeg" // Fallback
	}

	return data, mime, nil
}
