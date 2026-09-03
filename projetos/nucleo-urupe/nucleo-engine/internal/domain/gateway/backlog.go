/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"context"
	"log"
	"time"

	"nucleo-engine/internal/data/sqlite"
)

type BacklogProcessor struct {
	repo     *sqlite.Repository
	executor func(ctx context.Context, ev TriggerEvent) error
}

func NewBacklogProcessor(repo *sqlite.Repository, executor func(ctx context.Context, ev TriggerEvent) error) *BacklogProcessor {
	return &BacklogProcessor{
		repo:     repo,
		executor: executor,
	}
}

func (p *BacklogProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.Process(ctx)
		}
	}
}

func (p *BacklogProcessor) Process(ctx context.Context) {
	pending, err := p.repo.GetPendingTurns("", 5)
	if err != nil {
		log.Printf("[BacklogProcessor] Failed to fetch pending turns: %v", err)
		return
	}

	if len(pending) > 0 {
		log.Printf("[BacklogProcessor] Attempting to clear %d pending turns...", len(pending))
	}

	for _, t := range pending {
		ev := TriggerEvent{
			ChannelID:    t.ChannelID,
			MessageID:    t.MessageID,
			MessageRowID: t.MessageRowID,
			Timestamp:    t.Timestamp,
			IsReactive:   t.IsReactive,
		}

		err := p.executor(ctx, ev)
		if err == nil {
			log.Printf("[BacklogProcessor] Successfully cleared turn %s", t.MessageID)
			p.repo.DeletePendingTurn(t.ID)
		} else {
			log.Printf("[BacklogProcessor] Retry failed for %s: %v", t.MessageID, err)
			p.repo.UpdatePendingTurnStatus(t.ID, "pending", err.Error())
			// Stop processing backlog for now if we hit another error (likely still overloaded)
			return
		}
		
		time.Sleep(2 * time.Second)
	}
}
