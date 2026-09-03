/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package llm

import (
	"context"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	config       Config
	repo         MetricLogger
	configReader ConfigReader
	geminiPool   []*GeminiProvider
	geminiIdx    int
	groq         *GroqProvider
	openrouter   *OpenRouterProvider
	sem          *DynamicSemaphore
	mu           sync.Mutex
}

func NewClient(config Config, repo MetricLogger, configReader ConfigReader) *Client {
	maxConcurrent := 4
	if configReader != nil {
		if val, err := strconv.Atoi(configReader.GetConfig("llm_max_concurrent", "4")); err == nil {
			maxConcurrent = val
		}
	}

	c := &Client{
		config:       config,
		repo:         repo,
		configReader: configReader,
		sem:          NewDynamicSemaphore(maxConcurrent),
	}

	c.initProviders()

	if configReader != nil {
		go c.watchConfig()
	}

	return c
}

func (c *Client) watchConfig() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if val, err := strconv.Atoi(c.configReader.GetConfig("llm_max_concurrent", "4")); err == nil {
			c.sem.Resize(val)
		}
	}
}

// Embed delegates to the underlying provider (OpenRouter/OpenAI for now)
func (c *Client) Embed(ctx context.Context, input string) ([]float32, error) {
	if c.openrouter != nil {
		return c.openrouter.Embed(ctx, input)
	}
	return nil, nil
}

func (c *Client) Complete(ctx context.Context, messages []Message, stop []string, opts RequestOptions) (Completion, error) {
	if opts.Model == "" && c.config.OpenRouterModel != "" {
		opts.Model = c.config.OpenRouterModel
	}
	if c.openrouter != nil {
		return c.openrouter.Complete(ctx, messages, stop, opts)
	}
	if len(c.geminiPool) > 0 {
		return c.geminiPool[0].Complete(ctx, messages, stop, opts)
	}
	return Completion{}, nil
}
