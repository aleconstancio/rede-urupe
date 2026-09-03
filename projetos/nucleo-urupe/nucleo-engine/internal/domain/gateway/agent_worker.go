/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"context"
	"log"

	"github.com/aleconstancio/talos/v2/engine/core/agent"
	"nucleo-engine/internal/data/sqlite"
)

// AgentWorker wraps the Talos Agent pipeline for event processing.
type AgentWorker struct {
	pipeline      *agent.Agent
	repo          *sqlite.Repository
	pulseDetector *MazePulseDetector
	backlogStore  *SQLiteBacklogStore
	broadcaster   Broadcaster
}

func NewAgentWorker(
	pipeline *agent.Agent,
	repo *sqlite.Repository,
	pulseDetector *MazePulseDetector,
	backlogStore *SQLiteBacklogStore,
	broadcaster Broadcaster,
) *AgentWorker {
	return &AgentWorker{
		pipeline:      pipeline,
		repo:          repo,
		pulseDetector: pulseDetector,
		backlogStore:  backlogStore,
		broadcaster:   broadcaster,
	}
}

func (w *AgentWorker) ProcessEvent(ctx context.Context, ev TriggerEvent) error {
	input, err := (&TriggerIngester{Repo: w.repo}).Ingest(ctx, ev)
	if err != nil {
		return err
	}

	should, eventType, reason, err := w.pulseDetector.ShouldProcess(ctx, input)
	if err != nil {
		log.Printf("[AgentWorker] Pulse check error: %v", err)
		return nil
	}
	if !should {
		log.Printf("[AgentWorker] Skipped: %s", reason)
		return nil
	}
	input.EventType = string(eventType)

	result, err := w.pipeline.Run(ctx, ev)
	if err != nil {
		log.Printf("[AgentWorker] Pipeline error: %v", err)
		w.backlogStore.SavePendingTurn(ctx, &agent.PendingTurn{
			ChannelID: ev.ChannelID,
			MessageID: ev.MessageID,
			Content:   "",
			Status:    "pending",
		})
		return err
	}

	if w.broadcaster != nil {
		w.broadcaster.Broadcast("feed")
	}

	_ = result
	return nil
}

func (w *AgentWorker) SetBroadcaster(b Broadcaster) {
	w.broadcaster = b
}
