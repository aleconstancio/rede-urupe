/*
 * Copyright (c) 2026 Talos V2 Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"

	"nucleo-engine/internal/pkg/timeutil"
)

// MessageAnnotation represents metacognitive enrichment of a single message.
type MessageAnnotation struct {
	MessageID            string    `json:"message_id"`
	ChannelID            string    `json:"channel_id"`
	EpisodeID            string    `json:"episode_id"`
	AuthorID             string    `json:"author_id"`
	TopicTags            []string  `json:"topic_tags"`
	StanceTags           []string  `json:"stance_tags"`
	StyleTags            []string  `json:"style_tags"`
	EvidenceType         string    `json:"evidence_type"`
	DurabilityScore      float64   `json:"durability_score"`
	RetrievalScore       float64   `json:"retrieval_score"`
	HumorScore           float64   `json:"humor_score"`
	SarcasmScore         float64   `json:"sarcasm_score"`
	ContradictsMessageID string    `json:"contradicts_message_id"`
	SupersedesMessageID  string    `json:"supersedes_message_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// SaveMessageAnnotation persists metacognitive enrichment for a message.
func (r *Repository) SaveMessageAnnotation(tx *sql.Tx, ann MessageAnnotation) error {
	exec := r.getExecutor(tx)
	ann.MessageID = strings.TrimSpace(ann.MessageID)
	ann.ChannelID = strings.TrimSpace(ann.ChannelID)
	ann.AuthorID = strings.TrimSpace(ann.AuthorID)
	if ann.MessageID == "" || ann.ChannelID == "" || ann.AuthorID == "" {
		return nil
	}

	now := timeutil.Now()
	topicTagsJSON, _ := json.Marshal(ann.TopicTags)
	stanceTagsJSON, _ := json.Marshal(ann.StanceTags)
	styleTagsJSON, _ := json.Marshal(ann.StyleTags)

	_, err := exec.Exec(`
		INSERT INTO message_annotations (
			message_id, channel_id, episode_id, author_id, topic_tags_json, stance_tags_json, style_tags_json,
			evidence_type, durability_score, retrieval_score, humor_score, sarcasm_score,
			contradicts_message_id, supersedes_message_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(message_id) DO UPDATE SET
			topic_tags_json = excluded.topic_tags_json,
			stance_tags_json = excluded.stance_tags_json,
			style_tags_json = excluded.style_tags_json,
			durability_score = excluded.durability_score,
			retrieval_score = excluded.retrieval_score,
			humor_score = excluded.humor_score,
			sarcasm_score = excluded.sarcasm_score,
			contradicts_message_id = excluded.contradicts_message_id,
			supersedes_message_id = excluded.supersedes_message_id,
			updated_at = excluded.updated_at
	`, ann.MessageID, ann.ChannelID, ann.EpisodeID, ann.AuthorID, string(topicTagsJSON), string(stanceTagsJSON), string(styleTagsJSON),
		ann.EvidenceType, ann.DurabilityScore, ann.RetrievalScore, ann.HumorScore, ann.SarcasmScore,
		ann.ContradictsMessageID, ann.SupersedesMessageID, now, now)
	return err
}

// GetMessageAnnotations retrieves annotations for a list of message IDs.
func (r *Repository) GetMessageAnnotations(messageIDs []string) ([]MessageAnnotation, error) {
	if len(messageIDs) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(messageIDs))
	args := make([]any, len(messageIDs))
	for i, id := range messageIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	rows, err := r.db.Query(`
		SELECT message_id, channel_id, episode_id, author_id, topic_tags_json, stance_tags_json, style_tags_json,
		       evidence_type, durability_score, retrieval_score, humor_score, sarcasm_score,
		       contradicts_message_id, supersedes_message_id, created_at, updated_at
		FROM message_annotations
		WHERE message_id IN (`+strings.Join(placeholders, ",")+`)
	`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MessageAnnotation
	for rows.Next() {
		var ann MessageAnnotation
		var topicJSON, stanceJSON, styleJSON string
		err := rows.Scan(
			&ann.MessageID, &ann.ChannelID, &ann.EpisodeID, &ann.AuthorID, &topicJSON, &stanceJSON, &styleJSON,
			&ann.EvidenceType, &ann.DurabilityScore, &ann.RetrievalScore, &ann.HumorScore, &ann.SarcasmScore,
			&ann.ContradictsMessageID, &ann.SupersedesMessageID, &ann.CreatedAt, &ann.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(topicJSON), &ann.TopicTags)
		json.Unmarshal([]byte(stanceJSON), &ann.StanceTags)
		json.Unmarshal([]byte(styleJSON), &ann.StyleTags)
		out = append(out, ann)
	}
	return out, nil
}

// ListMessageAnnotations returns the most recent annotations for a channel.
func (r *Repository) ListMessageAnnotations(channelID string, limit int) ([]MessageAnnotation, error) {
	rows, err := r.db.Query(`
		SELECT message_id, channel_id, episode_id, author_id, topic_tags_json, stance_tags_json, style_tags_json,
		       evidence_type, durability_score, retrieval_score, humor_score, sarcasm_score,
		       contradicts_message_id, supersedes_message_id, created_at, updated_at
		FROM message_annotations
		WHERE channel_id = ?
		ORDER BY created_at DESC
		LIMIT ?
	`, channelID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []MessageAnnotation
	for rows.Next() {
		var ann MessageAnnotation
		var topicJSON, stanceJSON, styleJSON string
		err := rows.Scan(
			&ann.MessageID, &ann.ChannelID, &ann.EpisodeID, &ann.AuthorID, &topicJSON, &stanceJSON, &styleJSON,
			&ann.EvidenceType, &ann.DurabilityScore, &ann.RetrievalScore, &ann.HumorScore, &ann.SarcasmScore,
			&ann.ContradictsMessageID, &ann.SupersedesMessageID, &ann.CreatedAt, &ann.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(topicJSON), &ann.TopicTags)
		json.Unmarshal([]byte(stanceJSON), &ann.StanceTags)
		json.Unmarshal([]byte(styleJSON), &ann.StyleTags)
		out = append(out, ann)
	}
	return out, nil
}
