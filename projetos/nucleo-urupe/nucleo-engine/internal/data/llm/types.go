package llm

import (
	"context"

	"nucleo-engine/internal/config"
)

type Config = config.Config

type MetricLogger interface {
	LogMetric(mType, subType string, value float64, latency int, metadata map[string]any) error
}

type ConfigReader interface {
	GetConfig(key, fallback string) string
}

type DynamicSemaphore struct {
	ch chan struct{}
}

func NewDynamicSemaphore(n int) *DynamicSemaphore {
	if n <= 0 {
		n = 4
	}
	return &DynamicSemaphore{ch: make(chan struct{}, n)}
}

func (s *DynamicSemaphore) Resize(n int) {
	if n <= 0 {
		return
	}
	s.ch = make(chan struct{}, n)
}

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type ContentPart struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
	Data     []byte `json:"data,omitempty"`
	ImageURL *struct {
		URL string `json:"url"`
	} `json:"image_url,omitempty"`
}

type RequestOptions struct {
	Model            string      `json:"model"`
	Temperature      float64     `json:"temperature"`
	HasTemperature   bool        `json:"-"`
	MaxTokens        int         `json:"max_tokens"`
	ResponseMimeType string      `json:"response_mime_type"`
	ResponseSchema   interface{} `json:"response_schema"`
	ServiceTier      string      `json:"service_tier"`
	Tools            interface{} `json:"tools"`
	Purpose          string      `json:"purpose,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Completion struct {
	Text      string      `json:"text"`
	Reasoning string      `json:"reasoning,omitempty"`
	Model     string      `json:"model"`
	Usage     Usage       `json:"usage"`
	ToolCalls interface{} `json:"tool_calls,omitempty"`
}

type ChatRequest struct {
	Model       string      `json:"model"`
	Messages    []Message   `json:"messages"`
	Temperature float64     `json:"temperature"`
	MaxTokens   int         `json:"max_tokens"`
	Stop        []string    `json:"stop"`
	Tools       interface{} `json:"tools"`
}

type ChatResponse struct {
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Content   string      `json:"content"`
			Reasoning string      `json:"reasoning"`
			ToolCalls interface{} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
	Usage Usage `json:"usage"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

type GroqProvider struct{}

func (g *GroqProvider) Complete(ctx context.Context, messages []Message, stop []string, opts RequestOptions) (Completion, error) {
	return Completion{}, nil
}

func (c *Client) initProviders() {
	c.geminiPool = make([]*GeminiProvider, 0)
	if c.config.GeminiAPIKey != "" {
		if p, err := NewGeminiProvider(context.Background(), c.config.GeminiAPIKey); err == nil {
			c.geminiPool = append(c.geminiPool, p)
		}
	}
	if c.config.OpenRouterAPIKey != "" {
		c.openrouter = NewOpenRouterProvider(c.config)
	}
}
