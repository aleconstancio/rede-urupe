/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/aleconstancio/talos/v2/domain"
	"github.com/aleconstancio/talos/v2/engine/core/agent"
	"github.com/aleconstancio/talos/v2/engine/core/prompt"
	minotaur "github.com/aleconstancio/talos/v2/behavior/persona"
	"nucleo-engine/internal/config"
	"nucleo-engine/internal/data/llm"
	"nucleo-engine/internal/data/sqlite"
)

// promptContextStore carries PromptContext between pipeline stages.
// Keyed by PipelineInput.TurnID.
var promptContextStore sync.Map

func storePromptContext(turnID string, pc PromptContext) {
	promptContextStore.Store(turnID, pc)
}

func loadPromptContext(turnID string) (PromptContext, bool) {
	v, ok := promptContextStore.Load(turnID)
	if !ok {
		return PromptContext{}, false
	}
	pc, ok := v.(PromptContext)
	return pc, ok
}

func deletePromptContext(turnID string) {
	promptContextStore.Delete(turnID)
}

// TriggerIngester converts Discord TriggerEvents into PipelineInput.
type TriggerIngester struct {
	Repo *sqlite.Repository
}

func (i *TriggerIngester) Ingest(ctx context.Context, raw any) (*agent.PipelineInput, error) {
	ev, ok := raw.(TriggerEvent)
	if !ok {
		return nil, fmt.Errorf("expected TriggerEvent, got %T", raw)
	}

	msg, err := i.Repo.GetMessageByDiscordID(ev.MessageID)
	if err != nil {
		return nil, fmt.Errorf("get message: %w", err)
	}

	eventType := "ambient"
	if ev.IsReactive {
		eventType = "reactive"
	}

	return &agent.PipelineInput{
		TurnID:        fmt.Sprintf("turn_%d", ev.MessageRowID),
		ParticipantID: msg.AuthorID,
		ChannelID:     ev.ChannelID,
		Content:       msg.Content,
		Metadata: map[string]string{
			"message_row_id": strconv.FormatInt(ev.MessageRowID, 10),
			"message_id":     ev.MessageID,
			"is_reactive":    strconv.FormatBool(ev.IsReactive),
			"event_type":     eventType,
			"timestamp":      ev.Timestamp.Format(time.RFC3339),
		},
		RawPayload: ev,
	}, nil
}

// MazeContextAssembler builds the 3-layer memory context.
type MazeContextAssembler struct {
	Repo      *sqlite.Repository
	Assembler *PayloadAssembler
}

func (a *MazeContextAssembler) Assemble(ctx context.Context, input *agent.PipelineInput) (*prompt.AssembledContext, error) {
	state, err := a.Repo.GetGatewayState(input.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("get gateway state: %w", err)
	}

	messageRowID, _ := strconv.ParseInt(input.Metadata["message_row_id"], 10, 64)

	working, err := a.Repo.GetMessagesByRowRange(input.ChannelID, state.LastPulseRowID, messageRowID)
	if err != nil {
		return nil, fmt.Errorf("get working messages: %w", err)
	}
	if len(working) > 40 {
		working = working[len(working)-40:]
	}

	timestamp, _ := time.Parse(time.RFC3339, input.Metadata["timestamp"])

	payload, err := a.Assembler.BuildPromptContext(input.ChannelID, working, timestamp)
	if err != nil {
		return nil, fmt.Errorf("build prompt context: %w", err)
	}

	// Store PromptContext for downstream stages.
	storePromptContext(input.TurnID, payload)

	return &prompt.AssembledContext{
		Preamble:    input.TurnID,
		AuditID:     input.TurnID,
		AssembledAt: time.Now(),
	}, nil
}

// MazePersonaSelector wraps the minotaur.Resolver.
type MazePersonaSelector struct {
	Resolver *minotaur.Resolver
}

func (s *MazePersonaSelector) Select(ctx context.Context, input *agent.PipelineInput) (*minotaur.ResolvedPersona, error) {
	reason := "respond"
	if input.Metadata["event_type"] == "ambient" {
		reason = "ambient"
	}
	return s.Resolver.Resolve(ctx, input.ParticipantID, domain.ResolveContext{ChannelID: input.ChannelID, Mode: reason})
}

// MazeGater runs the ambient gate LLM and converts to GateResult.
type MazeGater struct {
	Gateway   *Gateway
	Assembler *PayloadAssembler
	GateModel string
}

func (g *MazeGater) ShouldProceed(ctx context.Context, input *agent.PipelineInput, assembledCtx *prompt.AssembledContext, persona *minotaur.ResolvedPersona) (*agent.GateResult, error) {
	// Reactive turns always proceed through the gater.
	if input.Metadata["event_type"] == "reactive" {
		return &agent.GateResult{ShouldProceed: true, Reason: "reactive"}, nil
	}

	payload, ok := loadPromptContext(input.TurnID)
	if !ok {
		return nil, fmt.Errorf("prompt context not found for turn %s", input.TurnID)
	}

	completion, err := g.Gateway.ExecuteStructured(
		ctx,
		g.GateModel,
		g.Assembler.BuildGateSystemPrompt(persona.Identity.Name),
		g.Assembler.BuildUserPromptPayload(ctx, input.Metadata["message_id"], g.Assembler.BuildGateUserPrompt(payload, time.Now())),
		GateSchema(),
		llm.RequestOptions{Purpose: "gate"},
	)
	if err != nil {
		return nil, fmt.Errorf("ambient gate failed: %w", err)
	}

	gateResp, err := ParseGateResponse(completion.Text)
	if err != nil {
		return nil, fmt.Errorf("parse gate response: %w", err)
	}

	return &agent.GateResult{
		ShouldProceed: gateResp.ShouldIntervene,
		Reason:        gateResp.ReasonCode,
		Metadata: map[string]any{
			"social_frame":      gateResp.SocialFrame,
			"invitation_signal": gateResp.InvitationSignal,
			"additive_value":    gateResp.AdditiveValue,
			"interruption_risk": gateResp.InterruptionRisk,
			"thinking_ledger":   gateResp.ThinkingLedger,
			"reason_code":       gateResp.ReasonCode,
			"confidence":        gateResp.Confidence,
			"gate_response":     gateResp,
		},
	}, nil
}

// MazeGenerator wraps the existing reply generation.
type MazeGenerator struct {
	Gateway   *Gateway
	Messenger *DiscordMessenger
	Assembler *PayloadAssembler
	Repo      *sqlite.Repository
	Cfg       config.Config
}

func (g *MazeGenerator) Generate(ctx context.Context, input *agent.PipelineInput, assembled *prompt.AssembledContext, persona *minotaur.ResolvedPersona, gateResult *agent.GateResult) (*agent.GatewayTurnResult, error) {
	payload, ok := loadPromptContext(input.TurnID)
	if !ok {
		return nil, fmt.Errorf("prompt context not found for turn %s", input.TurnID)
	}
	defer deletePromptContext(input.TurnID)

	isReactive := input.Metadata["event_type"] == "reactive"

	// Show typing.
	if g.Messenger.session != nil {
		_ = g.Messenger.session.ChannelTyping(input.ChannelID)
	}

	// Sync Discord profile.
	g.Messenger.SyncDiscordProfile(ctx, input.ChannelID, persona.Identity)

	// Build system prompt.
	systemPrompt := g.Assembler.AssembleReplyPrompt(*persona, isReactive)

	// Build gate context string for user prompt.
	var gateResp *GateResponse
	if gateResult != nil && gateResult.Metadata != nil {
		if gr, ok := gateResult.Metadata["gate_response"].(*GateResponse); ok {
			gateResp = gr
		}
	}

	// Build user prompt.
	timestamp, _ := time.Parse(time.RFC3339, input.Metadata["timestamp"])
	basePrompt := g.Assembler.BuildReplyUserPrompt(payload, isReactive, gateResp, timestamp)
	userPrompt := g.Assembler.BuildUserPromptPayload(ctx, input.Metadata["message_id"], basePrompt)

	// Execute LLM.
	completion, err := g.Gateway.ExecuteStructured(
		ctx,
		g.Cfg.GatewayReplyModel,
		systemPrompt,
		userPrompt,
		ReplySchema(),
		llm.RequestOptions{Purpose: "response"},
	)
	if err != nil {
		return nil, fmt.Errorf("LLM execution: %w", err)
	}

	// Parse response.
	resp, err := ParseReplyResponse(completion.Text)
	if err != nil {
		return nil, fmt.Errorf("parse reply: %w", err)
	}

	monologueJSON, _ := json.Marshal(resp.InternalMonologue)
	ledgerJSON, _ := json.Marshal(resp.GroundingLedger)

	// Send reply.
	err = g.Messenger.SendAndStoreReply(
		input.ChannelID,
		input.Metadata["message_id"],
		resp.ReplyText,
		string(monologueJSON),
		string(ledgerJSON),
		resp.SuggestedReactions,
	)
	if err != nil {
		return nil, fmt.Errorf("send reply: %w", err)
	}

	return &agent.GatewayTurnResult{
		ReplyText:          resp.ReplyText,
		SuggestedReactions: resp.SuggestedReactions,
		InternalMonologue:  string(monologueJSON),
		GroundingLedger:    resp.GroundingLedger,
		StanceUpdates:      resp.StanceUpdates,
	}, nil
}
