/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"nucleo-engine/internal/data/llm"
	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/pkg/timeutil"
)

type CapsuleWorker struct {
	repo      *sqlite.Repository
	llm       *llm.Client
	model     string
	channelID string
}

func NewCapsuleWorker(repo *sqlite.Repository, llmClient *llm.Client, model, channelID string) *CapsuleWorker {
	return &CapsuleWorker{
		repo:      repo,
		llm:       llmClient,
		model:     model,
		channelID: channelID,
	}
}

func (w *CapsuleWorker) Start(ctx context.Context) {
	log.Println("[CapsuleWorker] Started adaptive worker")
	w.runOnce(ctx)

	timer := time.NewTimer(w.nextRunDelay())
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			w.runOnce(ctx)
			timer.Reset(w.nextRunDelay())
		}
	}
}

func (w *CapsuleWorker) runOnce(ctx context.Context) {
	now := timeutil.Now()
	cutoff := now.Add(-7 * 24 * time.Hour)
	latest, err := w.repo.GetLatestUnassignedMessageByChannelSince(w.channelID, cutoff)
	if err != nil {
		log.Printf("[CapsuleWorker] Error fetching latest unassigned message: %v", err)
		return
	}
	if latest == nil {
		return
	}

	cfg := defaultCapsuleSegmenterConfig()
	if w.isPassiveMode() {
		cfg.MaxSegmentsPerRun = 1
	}

	messages, err := w.loadContiguousUnassignedTail(latest, cutoff, cfg)
	if err != nil {
		log.Printf("[CapsuleWorker] Error loading unassigned tail: %v", err)
		return
	}
	if len(messages) == 0 {
		return
	}

	segments := planCapsuleSegments(messages, now, cfg)
	if len(segments) == 0 {
		return
	}

	for _, segment := range segments {
		if ctx.Err() != nil {
			return
		}

		log.Printf(
			"[CapsuleWorker] Generating capsule for rows %d to %d (%d messages, reason=%s)",
			segment.StartExclusiveRowID+1,
			segment.EndInclusiveRowID,
			segment.MessageCount,
			segment.Reason,
		)

		capsule, usage, err := w.generateCapsule(ctx, segment)
		if err != nil {
			log.Printf("[CapsuleWorker] Capsule generation failed for rows %d-%d: %v", segment.StartExclusiveRowID+1, segment.EndInclusiveRowID, err)
			return
		}

		if err := w.commitSegment(segment, capsule); err != nil {
			log.Printf("[CapsuleWorker] Commit failed for rows %d-%d: %v", segment.StartExclusiveRowID+1, segment.EndInclusiveRowID, err)
			return
		}

		log.Printf("[CapsuleWorker] Saved capsule ID %d for rows %d-%d", capsule.ID, segment.StartExclusiveRowID+1, segment.EndInclusiveRowID)
		_ = w.repo.LogBudgetEvent(
			fmt.Sprintf("%s:%s", capsule.DayDate, capsule.TimeSpan),
			"capsule_generation",
			w.model,
			"memory",
			usage.PromptTokens,
			usage.CompletionTokens,
		)
	}
}

func (w *CapsuleWorker) isPassiveMode() bool {
	return w.repo.GetConfig("passive_mode", "true") == "true"
}

func (w *CapsuleWorker) nextRunDelay() time.Duration {
	if w.isPassiveMode() {
		return 5 * time.Minute
	}
	return nextHourlyRunDelay(timeutil.Now())
}

func (w *CapsuleWorker) loadContiguousUnassignedTail(latest *sqlite.Message, cutoff time.Time, cfg capsuleSegmenterConfig) ([]sqlite.Message, error) {
	if latest == nil {
		return nil, nil
	}

	probeLimit := cfg.TailProbeMessages
	minProbe := cfg.MaxMessages*cfg.MaxSegmentsPerRun + cfg.OrderingProbeWindow
	if minProbe < cfg.MaxMessages*2 {
		minProbe = cfg.MaxMessages * 2
	}
	if probeLimit < minProbe {
		probeLimit = minProbe
	}

	window, err := w.repo.GetRecentMessagesByChannelBeforeRowSince(w.channelID, latest.ID, cutoff, probeLimit)
	if err != nil {
		return nil, err
	}

	tailDesc := make([]sqlite.Message, 0, len(window))
	for _, msg := range window {
		if strings.TrimSpace(msg.EpisodeID) != "" {
			if len(tailDesc) > 0 {
				break
			}
			continue
		}
		tailDesc = append(tailDesc, msg)
	}
	if len(tailDesc) == 0 {
		return nil, nil
	}

	tail := make([]sqlite.Message, 0, len(tailDesc))
	for i := len(tailDesc) - 1; i >= 0; i-- {
		tail = append(tail, tailDesc[i])
	}
	return tail, nil
}

func (w *CapsuleWorker) generateCapsule(ctx context.Context, segment capsuleSegment) (sqlite.MemoryCapsule, llm.Usage, error) {
	prompt := buildCapsulePrompt(segment)
	schema := capsuleSchema()

	completion, err := w.llm.Complete(ctx, []llm.Message{{Role: "user", Content: prompt}}, nil, llm.RequestOptions{
		Model:            w.model,
		ResponseMimeType: "application/json",
		ResponseSchema:   schema,
		Purpose:          "summary",
	})
	if err != nil {
		return sqlite.MemoryCapsule{}, llm.Usage{}, err
	}

	var capsule sqlite.MemoryCapsule
	if err := json.Unmarshal([]byte(completion.Text), &capsule); err != nil {
		return sqlite.MemoryCapsule{}, llm.Usage{}, err
	}

	hydrateCapsuleFromSegment(&capsule, segment)
	return capsule, completion.Usage, nil
}

func (w *CapsuleWorker) commitSegment(segment capsuleSegment, capsule sqlite.MemoryCapsule) error {
	tx, err := w.repo.BeginImmediate()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	assignedCount, err := w.repo.CountAssignedMessagesByRowRangeTx(tx, w.channelID, segment.StartExclusiveRowID, segment.EndInclusiveRowID)
	if err != nil {
		return err
	}
	if assignedCount > 0 {
		return fmt.Errorf("segment rows %d-%d already contain %d assigned messages", segment.StartExclusiveRowID+1, segment.EndInclusiveRowID, assignedCount)
	}

	if err := w.repo.SaveMemoryCapsuleTx(tx, &capsule); err != nil {
		return err
	}
	if err := w.repo.AssignEpisodeToMessagesByRowRangeTx(tx, w.channelID, segment.StartExclusiveRowID, segment.EndInclusiveRowID, capsule.EpisodeID); err != nil {
		return err
	}

	return tx.Commit()
}
