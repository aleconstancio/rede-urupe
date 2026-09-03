package gateway

import (
	"context"

	talosgate "github.com/aleconstancio/talos/v2/engine/core/gateway"
	"nucleo-engine/internal/data/llm"
	"google.golang.org/genai"
)

// Type aliases — replaced by Talos gateway library types.
type GateResponse = talosgate.GateResponse
type ReplyInternalMonologue = talosgate.ReplyInternalMonologue
type ReplyResponse = talosgate.ReplyResponse

var ParseGateResponse = talosgate.ParseGateResponse
var ParseReplyResponse = talosgate.ParseReplyResponse

type Gateway struct {
	llm *llm.Client
}

func NewGateway(llmClient *llm.Client) *Gateway {
	return &Gateway{llm: llmClient}
}

func (g *Gateway) ExecuteStructured(ctx context.Context, model string, systemPrompt string, userPrompt interface{}, schema *genai.Schema, opts ...llm.RequestOptions) (llm.Completion, error) {
	var requestOpts llm.RequestOptions
	if len(opts) > 0 {
		requestOpts = opts[0]
	}
	requestOpts.Model = model
	requestOpts.ResponseMimeType = "application/json"
	requestOpts.ResponseSchema = schema

	messages := []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	return g.llm.Complete(ctx, messages, nil, requestOpts)
}
