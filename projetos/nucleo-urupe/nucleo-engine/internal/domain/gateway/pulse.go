/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"database/sql"
	"errors"
	"time"

	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/pkg/timeutil"
)

func (w *GatewayWorker) seedState() error {
	state, err := w.repo.GetGatewayState(w.channelID)
	if err == nil {
		w.state = state
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	latest, err := w.repo.GetLatestMessageByChannel(w.channelID)
	if err != nil {
		return err
	}

	now := timeutil.Now()
	w.state = &sqlite.GatewayState{
		ChannelID:      w.channelID,
		LastPulseRowID: 0,
		LastPulseMsgID: "",
		LastPulseAt:    now,
		UpdatedAt:      now,
	}
	if latest != nil {
		w.state.LastPulseRowID = latest.ID
		w.state.LastPulseMsgID = latest.DiscordID
		w.state.LastPulseAt = timeutil.InBrasilia(latest.Timestamp)
	}

	return w.repo.UpsertGatewayState(w.state)
}

func (w *GatewayWorker) persistPulse(ev TriggerEvent) error {
	w.state = &sqlite.GatewayState{
		ChannelID:      w.channelID,
		LastPulseRowID: ev.MessageRowID,
		LastPulseMsgID: ev.MessageID,
		LastPulseAt:    timeutil.InBrasilia(ev.Timestamp),
		UpdatedAt:      timeutil.Now(),
	}
	return w.repo.UpsertGatewayState(w.state)
}

func (w *GatewayWorker) shouldTriggerAmbient(ev TriggerEvent) bool {
	if ev.IsReactive {
		return false
	}
	
	isPassive := w.repo.GetConfig("passive_mode", "true") == "true"
	if isPassive {
		return false
	}

	// Heuristics: 20 messages since last pulse OR 1 hour passed
	// But must be at least 3 minutes since last pulse to avoid spamming
	return (w.messagesSinceLastPulse >= 20 || ev.Timestamp.After(w.state.LastPulseAt.Add(1*time.Hour))) &&
		ev.Timestamp.After(w.state.LastPulseAt.Add(3*time.Minute))
}
