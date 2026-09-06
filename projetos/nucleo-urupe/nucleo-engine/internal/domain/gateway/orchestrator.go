/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aleconstancio/talos/v2/domain"
	"nucleo-engine/internal/data/llm"
)

func (w *GatewayWorker) runBacklogTurn(ctx context.Context, ev TriggerEvent) error {
	working, err := w.repo.GetMessagesByRowRange(ev.ChannelID, w.state.LastPulseRowID, ev.MessageRowID)
	if err != nil {
		return err
	}
	payload, err := w.assembler.BuildPromptContext(ev.ChannelID, working, ev.Timestamp)
	if err != nil {
		return err
	}

	if ev.IsReactive {
		return w.runReactiveTurn(ctx, ev, payload, true)
	}
	
	gateResp, err := w.runAmbientGate(ctx, ev, payload)
	if err == nil && gateResp.ShouldIntervene {
		return w.runAmbientReply(ctx, ev, payload, gateResp)
	}
	return err
}

func (w *GatewayWorker) runReactiveTurn(ctx context.Context, ev TriggerEvent, payload PromptContext, isRetry bool) error {
	channelID := w.getChannelID(ev)
	if w.messenger.session != nil {
		_ = w.messenger.session.ChannelTyping(channelID)
	}
	
	resolved, err := w.resolver.Resolve(ctx, channelID, domain.ResolveContext{ChannelID: channelID, Mode: "respond"}) 
	if err != nil {
		log.Printf("[GatewayWorker] Resolver error: %v", err)
	}

	w.messenger.SyncDiscordProfile(ctx, channelID, resolved.Identity)

	completion, err := w.gateway.ExecuteStructured(
		ctx,
		w.replyModel,
		w.assembler.AssembleReplyPrompt(*resolved, true),
		w.assembler.BuildUserPromptPayload(ctx, ev.MessageID, w.assembler.BuildReplyUserPrompt(payload, true, nil, ev.Timestamp)),
		ReplySchema(),
		llm.RequestOptions{Purpose: "response"},
	)
	if err != nil {
		return w.handleReplyError(ev, completion, err, isRetry)
	}
	
	_ = w.logBudget(ev, "gateway_reply", "reactive", w.replyModel, completion)

	resp, err := ParseReplyResponse(completion.Text)
	if err != nil {
		return err
	}
	if !resp.ShouldIntervene || (strings.TrimSpace(resp.ReplyText) == "" && len(resp.SuggestedReactions) == 0) {
		log.Printf("[GatewayWorker] Reactive turn produced no reply/reactions for %s", ev.MessageID)
		return nil
	}

	return w.sendReply(ev, resp)
}

func (w *GatewayWorker) runAmbientGate(ctx context.Context, ev TriggerEvent, payload PromptContext) (*GateResponse, error) {
	channelID := w.getChannelID(ev)
	resolved, err := w.resolver.Resolve(ctx, channelID, domain.ResolveContext{ChannelID: channelID, Mode: "gate"})
	if err != nil {
		log.Printf("[GatewayWorker] Resolver error in gate: %v", err)
	}

	completion, err := w.gateway.ExecuteStructured(
		ctx,
		w.gateModel,
		w.assembler.BuildGateSystemPrompt(resolved.Identity.Name),
		w.assembler.BuildUserPromptPayload(ctx, ev.MessageID, w.assembler.BuildGateUserPrompt(payload, ev.Timestamp)),
		GateSchema(),
		llm.RequestOptions{Purpose: "response"}, // Gating is similar to response triage
	)
	if err != nil {
		return nil, fmt.Errorf("ambient gate failed: %w", err)
	}
	
	_ = w.logBudget(ev, "gateway_gate", "ambient", w.gateModel, completion)

	return ParseGateResponse(completion.Text)
}

func (w *GatewayWorker) runAmbientReply(ctx context.Context, ev TriggerEvent, payload PromptContext, gateResp *GateResponse) error {
	channelID := w.getChannelID(ev)
	if w.messenger.session != nil {
		_ = w.messenger.session.ChannelTyping(channelID)
	}

	resolved, err := w.resolver.Resolve(ctx, channelID, domain.ResolveContext{ChannelID: channelID, Mode: gateResp.ReasonCode})
	if err != nil {
		log.Printf("[GatewayWorker] Resolver error: %v", err)
	}

	w.messenger.SyncDiscordProfile(ctx, channelID, resolved.Identity)

	completion, err := w.gateway.ExecuteStructured(
		ctx,
		w.replyModel,
		w.assembler.AssembleReplyPrompt(*resolved, false),
		w.assembler.BuildUserPromptPayload(ctx, ev.MessageID, w.assembler.BuildReplyUserPrompt(payload, false, gateResp, ev.Timestamp)),
		ReplySchema(),
		llm.RequestOptions{Purpose: "response"},
	)
	if err != nil {
		return fmt.Errorf("ambient reply failed: %w", err)
	}
	
	_ = w.logBudget(ev, "gateway_reply", "ambient", w.replyModel, completion)

	resp, err := ParseReplyResponse(completion.Text)
	if err != nil {
		return err
	}
	if !resp.ShouldIntervene || (strings.TrimSpace(resp.ReplyText) == "" && len(resp.SuggestedReactions) == 0) {
		log.Printf("[GatewayWorker] Ambient reply produced no reply/reactions for %s (reason=%s)", ev.MessageID, gateResp.ReasonCode)
		return nil
	}

	return w.sendReply(ev, resp)
}

func (w *GatewayWorker) handleReplyError(ev TriggerEvent, completion llm.Completion, err error, isRetry bool) error {
	channelID := w.getChannelID(ev)
	errMsg := strings.ToLower(err.Error())
	if strings.Contains(errMsg, "503") || strings.Contains(errMsg, "overloaded") || strings.Contains(errMsg, "high demand") {
		log.Printf("[GatewayWorker] Deferring turn due to API overload (%s): %v", completion.Model, err)
		
		if saveErr := w.repo.SavePendingTurn(channelID, ev.MessageID, ev.MessageRowID, ev.Timestamp, ev.IsReactive, err.Error()); saveErr != nil {
			log.Printf("[GatewayWorker] Failed to save pending turn: %v", saveErr)
		}

		if !isRetry {
			model := completion.Model
			if model == "" {
				model = w.replyModel
			}
			msg := fmt.Sprintf("Minha API(%s) está sobrecarregada, te respondo quando meu sinal voltar", model)
			_ = w.messenger.SendAndStoreReply(channelID, ev.MessageID, msg, "API Overload Deferral", "{}", nil)
		}
		return nil 
	}
	return fmt.Errorf("reply failed: %w", err)
}

func (w *GatewayWorker) sendReply(ev TriggerEvent, resp *ReplyResponse) error {
	channelID := w.getChannelID(ev)
	monologueJSON, _ := json.Marshal(resp.InternalMonologue)
	ledgerJSON, _ := json.Marshal(resp.GroundingLedger)
	return w.messenger.SendAndStoreReply(channelID, ev.MessageID, resp.ReplyText, string(monologueJSON), string(ledgerJSON), resp.SuggestedReactions)
}

func (w *GatewayWorker) logBudget(ev TriggerEvent, reason, triggerType, model string, completion llm.Completion) error {
	return w.repo.LogBudgetEvent(fmt.Sprintf("%s:%s", ev.MessageID, reason), reason, model, triggerType, completion.Usage.PromptTokens, completion.Usage.CompletionTokens)
}
