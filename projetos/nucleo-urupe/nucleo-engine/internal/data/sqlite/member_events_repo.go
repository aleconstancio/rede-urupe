/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"time"
)

type MemberEvent struct {
	ID        int64  `json:"id"`
	GuildID   string `json:"guild_id"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	EventType string `json:"event_type"` // join, leave, kick, ban, unban
	Reason    string `json:"reason"`
	ExecutorID string `json:"executor_id"`
	CreatedAt string `json:"created_at"`
}

func (r *Repository) LogMemberEvent(guildID, userID, userName, eventType, reason, executorID string) error {
	_, err := r.db.Exec(`
		INSERT INTO member_events (guild_id, user_id, user_name, event_type, reason, executor_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`, guildID, userID, userName, eventType, reason, executorID)
	return err
}

func (r *Repository) GetMemberEvents(guildID string, limit int) ([]MemberEvent, error) {
	rows, err := r.db.Query(`
		SELECT id, guild_id, user_id, user_name, event_type, reason, executor_id, created_at
		FROM member_events
		WHERE guild_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, guildID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []MemberEvent
	for rows.Next() {
		var e MemberEvent
		if err := rows.Scan(&e.ID, &e.GuildID, &e.UserID, &e.UserName, &e.EventType, &e.Reason, &e.ExecutorID, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *Repository) GetMemberEventsByUser(guildID, userID string, limit int) ([]MemberEvent, error) {
	rows, err := r.db.Query(`
		SELECT id, guild_id, user_id, user_name, event_type, reason, executor_id, created_at
		FROM member_events
		WHERE guild_id = ? AND user_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, guildID, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []MemberEvent
	for rows.Next() {
		var e MemberEvent
		if err := rows.Scan(&e.ID, &e.GuildID, &e.UserID, &e.UserName, &e.EventType, &e.Reason, &e.ExecutorID, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *Repository) GetMemberEventsByType(guildID, eventType string, limit int) ([]MemberEvent, error) {
	rows, err := r.db.Query(`
		SELECT id, guild_id, user_id, user_name, event_type, reason, executor_id, created_at
		FROM member_events
		WHERE guild_id = ? AND event_type = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, guildID, eventType, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []MemberEvent
	for rows.Next() {
		var e MemberEvent
		if err := rows.Scan(&e.ID, &e.GuildID, &e.UserID, &e.UserName, &e.EventType, &e.Reason, &e.ExecutorID, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *Repository) GetMemberGrowthDaily(guildID string, days int) ([]DailyGrowth, error) {
	since := time.Now().AddDate(0, 0, -days).UTC().Format("2006-01-02")
	rows, err := r.db.Query(`
		SELECT DATE(created_at) as day,
			SUM(CASE WHEN event_type = 'join' THEN 1 ELSE 0 END) as joins,
			SUM(CASE WHEN event_type IN ('leave', 'kick', 'ban') THEN 1 ELSE 0 END) as leaves
		FROM member_events
		WHERE guild_id = ? AND DATE(created_at) >= ?
		GROUP BY DATE(created_at)
		ORDER BY day ASC
	`, guildID, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var growth []DailyGrowth
	for rows.Next() {
		var d DailyGrowth
		if err := rows.Scan(&d.Date, &d.Joins, &d.Leaves); err != nil {
			return nil, err
		}
		d.NetChange = d.Joins - d.Leaves
		growth = append(growth, d)
	}
	return growth, nil
}

type DailyGrowth struct {
	Date      string `json:"date"`
	Joins     int    `json:"joins"`
	Leaves    int    `json:"leaves"`
	NetChange int    `json:"net_change"`
}
