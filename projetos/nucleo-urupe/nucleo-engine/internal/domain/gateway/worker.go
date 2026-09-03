/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"

	minotaur "github.com/aleconstancio/talos/v2/behavior/persona"
	"nucleo-engine/internal/data/sqlite"
)

type Broadcaster interface {
	Broadcast(event string)
}

type TriggerEvent struct {
	ChannelID    string
	MessageRowID int64
	MessageID    string
	Timestamp    time.Time
	IsReactive   bool
}

type GatewayWorker struct {
	repo                   *sqlite.Repository
	gateway                *Gateway
	assembler              *PayloadAssembler
	messenger              *DiscordMessenger
	backlog                *BacklogProcessor
	TriggerChan            chan TriggerEvent
	channelID              string
	guildID                string
	gateModel              string
	replyModel             string
	messagesSinceLastPulse int
	state                  *sqlite.GatewayState
	resolver               *minotaur.Resolver
	activeChannels         map[string]bool
}

func NewGatewayWorker(repo *sqlite.Repository, gateway *Gateway, guildID, channelID, gateModel, replyModel string) *GatewayWorker {
	w := &GatewayWorker{
		repo:           repo,
		gateway:        gateway,
		TriggerChan:    make(chan TriggerEvent, 100),
		guildID:        guildID,
		channelID:      channelID,
		gateModel:      gateModel,
		replyModel:     replyModel,
		resolver:       minotaur.NewResolver(repo),
		assembler:      NewPayloadAssembler(repo),
		activeChannels: map[string]bool{channelID: true},
	}

	w.messenger = NewDiscordMessenger(nil, repo, nil)
	w.backlog = NewBacklogProcessor(repo, w.runBacklogTurn)

	return w
}

func (w *GatewayWorker) SetSession(s *discordgo.Session) {
	w.messenger.session = s
}

func (w *GatewayWorker) SetBroadcaster(b Broadcaster) {
	w.messenger.broadcaster = b
}

func (w *GatewayWorker) Start(ctx context.Context) {
	if err := w.seedState(); err != nil {
		log.Printf("[GatewayWorker] Failed to seed gateway state: %v", err)
	}

	// Initial emoji sync
	w.SyncEmojis(ctx)

	log.Println("[GatewayWorker] Started sequential worker")

	go w.backlog.Start(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-w.TriggerChan:
			if err := w.processEvent(ctx, ev); err != nil {
				log.Printf("[GatewayWorker] Event error: %v", err)
			}
		}
	}
}

func (w *GatewayWorker) processEvent(ctx context.Context, ev TriggerEvent) error {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	if !w.isChannelActive(ev.ChannelID) {
		return nil
	}
	if w.state == nil {
		if err := w.seedState(); err != nil {
			return err
		}
	}

	w.messagesSinceLastPulse++
	triggerAmbient := w.shouldTriggerAmbient(ev)

	if !ev.IsReactive && !triggerAmbient {
		w.messenger.Broadcast("feed")
		return nil
	}

	// Load memory context
	working, err := w.repo.GetMessagesByRowRange(w.channelID, w.state.LastPulseRowID, ev.MessageRowID)
	if err != nil {
		return fmt.Errorf("load working memory: %w", err)
	}
	if len(working) > 40 {
		working = working[len(working)-40:]
	}

	payload, err := w.assembler.BuildPromptContext(w.channelID, working, ev.Timestamp)
	if err != nil {
		return err
	}

	turnCompleted := false
	switch {
	case ev.IsReactive:
		if err := w.runReactiveTurn(ctx, ev, payload, false); err != nil {
			_ = w.messenger.SendAndStoreReply(w.channelID, ev.MessageID, "minha api não respondeu depois eu respondo :bcatdespair:", "", "", nil)
			return err
		}
		turnCompleted = true
	case triggerAmbient:
		gateResp, err := w.runAmbientGate(ctx, ev, payload)
		if err != nil {
			// Don't spam ambient errors, but log them
			log.Printf("[GatewayWorker] Ambient Gate failed: %v", err)
			return err
		}
		if gateResp.ShouldIntervene {
			if err := w.runAmbientReply(ctx, ev, payload, gateResp); err != nil {
				_ = w.messenger.SendAndStoreReply(w.channelID, ev.MessageID, "minha api não respondeu depois eu respondo :bcatdespair:", "", "", nil)
				return err
			}
		}
		turnCompleted = true
	}

	if turnCompleted {
		if err := w.persistPulse(ev); err != nil {
			return err
		}
		w.messagesSinceLastPulse = 0
		w.messenger.Broadcast("metrics")
	}

	w.messenger.Broadcast("feed")
	return nil
}

func (w *GatewayWorker) isChannelActive(channelID string) bool {
	return w.activeChannels[channelID]
}

func (w *GatewayWorker) AddChannel(channelID string) {
	if w.activeChannels == nil {
		w.activeChannels = make(map[string]bool)
	}
	w.activeChannels[channelID] = true
}

func (w *GatewayWorker) SyncEmojis(ctx context.Context) {
	if w.messenger.session == nil || w.guildID == "" {
		return
	}

	emojis, err := w.messenger.session.GuildEmojis(w.guildID)
	if err != nil {
		log.Printf("[GatewayWorker] Failed to sync emojis: %v", err)
		return
	}

	emojiMap := make(map[string]string)
	for _, e := range emojis {
		emojiMap[e.Name] = e.ID
	}

	w.assembler.SetEmojis(emojiMap)
	w.messenger.SetEmojis(emojiMap)
	log.Printf("[GatewayWorker] Synced %d custom emojis from guild %s", len(emojiMap), w.guildID)
}
