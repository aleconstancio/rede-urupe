/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"encoding/json"
)

type AuditEvent struct {
	ID        int64  `json:"id"`
	GuildID   string `json:"guild_id"`
	ActionType string `json:"action_type"` // moderation, member, config, role
	Action    string `json:"action"`       // warn, delete, timeout, kick, ban, join, leave, config_change, role_assign, role_remove
	ActorID   string `json:"actor_id"`
	ActorName string `json:"actor_name"`
	TargetID  string `json:"target_id"`
	TargetName string `json:"target_name"`
	Details   string `json:"details"`
	ChannelID string `json:"channel_id"`
	CreatedAt string `json:"created_at"`
}

func (r *Repository) LogAuditEvent(guildID, actionType, action, actorID, actorName, targetID, targetName, details, channelID string) error {
	_, err := r.db.Exec(`
		INSERT INTO audit_log (guild_id, action_type, action, actor_id, actor_name, target_id, target_name, details, channel_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, guildID, actionType, action, actorID, actorName, targetID, targetName, details, channelID)
	return err
}

func (r *Repository) GetAuditLog(guildID string, limit int) ([]AuditEvent, error) {
	rows, err := r.db.Query(`
		SELECT id, guild_id, action_type, action, actor_id, actor_name, target_id, target_name, details, channel_id, created_at
		FROM audit_log
		WHERE guild_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, guildID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAuditEvents(rows)
}

func (r *Repository) GetAuditLogByUser(guildID, userID string, limit int) ([]AuditEvent, error) {
	rows, err := r.db.Query(`
		SELECT id, guild_id, action_type, action, actor_id, actor_name, target_id, target_name, details, channel_id, created_at
		FROM audit_log
		WHERE guild_id = ? AND (actor_id = ? OR target_id = ?)
		ORDER BY created_at DESC
		LIMIT ?
	`, guildID, userID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAuditEvents(rows)
}

func (r *Repository) GetAuditLogByAction(guildID, actionType string, limit int) ([]AuditEvent, error) {
	rows, err := r.db.Query(`
		SELECT id, guild_id, action_type, action, actor_id, actor_name, target_id, target_name, details, channel_id, created_at
		FROM audit_log
		WHERE guild_id = ? AND action_type = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, guildID, actionType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAuditEvents(rows)
}

func (r *Repository) GetWarningsByUser(guildID, userID string) ([]WarningEvent, error) {
	rows, err := r.db.Query(`
		SELECT id, guild_id, user_id, channel_id, reason, message_id, severity, created_at
		FROM guild_warnings
		WHERE guild_id = ? AND user_id = ?
		ORDER BY created_at DESC
	`, guildID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warnings []WarningEvent
	for rows.Next() {
		var w WarningEvent
		if err := rows.Scan(&w.ID, &w.GuildID, &w.UserID, &w.ChannelID, &w.Reason, &w.MessageID, &w.Severity, &w.CreatedAt); err != nil {
			return nil, err
		}
		warnings = append(warnings, w)
	}
	return warnings, nil
}

type WarningEvent struct {
	ID        int64  `json:"id"`
	GuildID   string `json:"guild_id"`
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	Reason    string `json:"reason"`
	MessageID string `json:"message_id"`
	Severity  int    `json:"severity"`
	CreatedAt string `json:"created_at"`
}

func (r *Repository) scanAuditEvents(rows interface{ Next() bool; Scan(...interface{}) error }) ([]AuditEvent, error) {
	var events []AuditEvent
	for rows.Next() {
		var e AuditEvent
		if err := rows.Scan(&e.ID, &e.GuildID, &e.ActionType, &e.Action, &e.ActorID, &e.ActorName, &e.TargetID, &e.TargetName, &e.Details, &e.ChannelID, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func AuditDetailsJSON(details map[string]interface{}) string {
	b, _ := json.Marshal(details)
	return string(b)
}
