/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */

package sqlite

import (
	"encoding/json"

	"nucleo-engine/internal/pkg/timeutil"
)

type ForumTemplate struct {
	ID               string            `json:"id"`
	GuildID          string            `json:"guild_id"`
	ChannelID        string            `json:"channel_id"`
	Title            string            `json:"title"`
	Body             string            `json:"body"`
	Tags             []string          `json:"tags"`
	Variables        []string          `json:"variables"`
	Schedule         string            `json:"schedule"`          // "manual", "daily", "weekly", "monthly"
	ScheduleConfig   map[string]any    `json:"schedule_config"`
	IsEnabled        bool              `json:"is_enabled"`
	CreatedAt        string            `json:"created_at"`
	UpdatedAt        string            `json:"updated_at"`
}

type ForumPost struct {
	ID               int64             `json:"id"`
	TemplateID       string            `json:"template_id"`
	GuildID          string            `json:"guild_id"`
	ChannelID        string            `json:"channel_id"`
	DiscordMessageID string            `json:"discord_message_id"`
	DiscordThreadID  string            `json:"discord_thread_id"`
	Title            string            `json:"title"`
	Body             string            `json:"body"`
	Tags             []string          `json:"tags"`
	Status           string            `json:"status"` // "draft", "published", "failed"
	Error            string            `json:"error"`
	CreatedAt        string            `json:"created_at"`
}

func (r *Repository) ensureForumSchema() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS forum_templates (
			id TEXT PRIMARY KEY,
			guild_id TEXT NOT NULL,
			channel_id TEXT NOT NULL,
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			tags_json TEXT NOT NULL DEFAULT '[]',
			variables_json TEXT NOT NULL DEFAULT '[]',
			schedule TEXT NOT NULL DEFAULT 'manual',
			schedule_config_json TEXT NOT NULL DEFAULT '{}',
			is_enabled BOOLEAN NOT NULL DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS forum_posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			template_id TEXT NOT NULL,
			guild_id TEXT NOT NULL,
			channel_id TEXT NOT NULL,
			discord_message_id TEXT NOT NULL DEFAULT '',
			discord_thread_id TEXT NOT NULL DEFAULT '',
			title TEXT NOT NULL,
			body TEXT NOT NULL,
			tags_json TEXT NOT NULL DEFAULT '[]',
			status TEXT NOT NULL DEFAULT 'draft',
			error TEXT NOT NULL DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(template_id) REFERENCES forum_templates(id)
		);
	`)
	return err
}

func (r *Repository) ListForumTemplates(guildID string) ([]ForumTemplate, error) {
	if err := r.ensureForumSchema(); err != nil {
		return nil, err
	}
	rows, err := r.db.Query(`
		SELECT id, guild_id, channel_id, title, body, tags_json, variables_json,
		       schedule, schedule_config_json, is_enabled, created_at, updated_at
		FROM forum_templates WHERE guild_id = ?
		ORDER BY created_at DESC
	`, guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []ForumTemplate
	for rows.Next() {
		var t ForumTemplate
		var tagsRaw, varsRaw, schedRaw string
		if err := rows.Scan(&t.ID, &t.GuildID, &t.ChannelID, &t.Title, &t.Body,
			&tagsRaw, &varsRaw, &t.Schedule, &schedRaw, &t.IsEnabled, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(tagsRaw), &t.Tags)
		json.Unmarshal([]byte(varsRaw), &t.Variables)
		json.Unmarshal([]byte(schedRaw), &t.ScheduleConfig)
		if t.ScheduleConfig == nil {
			t.ScheduleConfig = make(map[string]any)
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *Repository) GetForumTemplate(id string) (*ForumTemplate, error) {
	if err := r.ensureForumSchema(); err != nil {
		return nil, err
	}
	t := &ForumTemplate{}
	var tagsRaw, varsRaw, schedRaw string
	err := r.db.QueryRow(`
		SELECT id, guild_id, channel_id, title, body, tags_json, variables_json,
		       schedule, schedule_config_json, is_enabled, created_at, updated_at
		FROM forum_templates WHERE id = ?
	`, id).Scan(&t.ID, &t.GuildID, &t.ChannelID, &t.Title, &t.Body,
		&tagsRaw, &varsRaw, &t.Schedule, &schedRaw, &t.IsEnabled, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(tagsRaw), &t.Tags)
	json.Unmarshal([]byte(varsRaw), &t.Variables)
	json.Unmarshal([]byte(schedRaw), &t.ScheduleConfig)
	if t.ScheduleConfig == nil {
		t.ScheduleConfig = make(map[string]any)
	}
	return t, nil
}

func (r *Repository) UpsertForumTemplate(t ForumTemplate) error {
	if err := r.ensureForumSchema(); err != nil {
		return err
	}
	now := timeutil.Now()
	tagsRaw, _ := json.Marshal(t.Tags)
	varsRaw, _ := json.Marshal(t.Variables)
	schedRaw, _ := json.Marshal(t.ScheduleConfig)

	_, err := r.db.Exec(`
		INSERT INTO forum_templates (id, guild_id, channel_id, title, body, tags_json, variables_json,
		                             schedule, schedule_config_json, is_enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			title = excluded.title, body = excluded.body,
			tags_json = excluded.tags_json, variables_json = excluded.variables_json,
			schedule = excluded.schedule, schedule_config_json = excluded.schedule_config_json,
			is_enabled = excluded.is_enabled, updated_at = excluded.updated_at
	`, t.ID, t.GuildID, t.ChannelID, t.Title, t.Body, tagsRaw, varsRaw,
		t.Schedule, schedRaw, t.IsEnabled, now, now)
	return err
}

func (r *Repository) DeleteForumTemplate(id string) error {
	if err := r.ensureForumSchema(); err != nil {
		return err
	}
	_, err := r.db.Exec("DELETE FROM forum_templates WHERE id = ?", id)
	return err
}

func (r *Repository) SaveForumPost(p ForumPost) (int64, error) {
	if err := r.ensureForumSchema(); err != nil {
		return 0, err
	}
	tagsRaw, _ := json.Marshal(p.Tags)
	res, err := r.db.Exec(`
		INSERT INTO forum_posts (template_id, guild_id, channel_id, discord_message_id, discord_thread_id,
		                         title, body, tags_json, status, error, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, p.TemplateID, p.GuildID, p.ChannelID, p.DiscordMessageID, p.DiscordThreadID,
		p.Title, p.Body, tagsRaw, p.Status, p.Error, timeutil.Now())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *Repository) ListForumPosts(guildID string, limit int) ([]ForumPost, error) {
	if err := r.ensureForumSchema(); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.Query(`
		SELECT id, template_id, guild_id, channel_id, discord_message_id, discord_thread_id,
		       title, body, tags_json, status, error, created_at
		FROM forum_posts WHERE guild_id = ?
		ORDER BY created_at DESC LIMIT ?
	`, guildID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []ForumPost
	for rows.Next() {
		var p ForumPost
		var tagsRaw string
		if err := rows.Scan(&p.ID, &p.TemplateID, &p.GuildID, &p.ChannelID,
			&p.DiscordMessageID, &p.DiscordThreadID, &p.Title, &p.Body,
			&tagsRaw, &p.Status, &p.Error, &p.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(tagsRaw), &p.Tags)
		posts = append(posts, p)
	}
	return posts, nil
}

func (r *Repository) GetScheduledForumTemplates() ([]ForumTemplate, error) {
	if err := r.ensureForumSchema(); err != nil {
		return nil, err
	}
	rows, err := r.db.Query(`
		SELECT id, guild_id, channel_id, title, body, tags_json, variables_json,
		       schedule, schedule_config_json, is_enabled, created_at, updated_at
		FROM forum_templates
		WHERE is_enabled = 1 AND schedule != 'manual'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []ForumTemplate
	for rows.Next() {
		var t ForumTemplate
		var tagsRaw, varsRaw, schedRaw string
		if err := rows.Scan(&t.ID, &t.GuildID, &t.ChannelID, &t.Title, &t.Body,
			&tagsRaw, &varsRaw, &t.Schedule, &schedRaw, &t.IsEnabled, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(tagsRaw), &t.Tags)
		json.Unmarshal([]byte(varsRaw), &t.Variables)
		json.Unmarshal([]byte(schedRaw), &t.ScheduleConfig)
		templates = append(templates, t)
	}
	return templates, nil
}
