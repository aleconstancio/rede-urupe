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

const (
	MemoryKindEpisodeSnapshot = "episode_snapshot"
	MemoryKindCaptureSlice    = "capture_slice"
	MemoryKindDailySummary    = "daily_summary"
	MemoryKindTopicSummary    = "topic_summary"
	MemoryKindCategorySummary = "category_summary"
)

func IsSummaryKind(kind string) bool {
	return kind == MemoryKindDailySummary || kind == MemoryKindTopicSummary || kind == MemoryKindCategorySummary
}

type TypedFacts struct {
	UserPreferences []string `json:"user_preferences,omitempty"`
	SystemEvents    []string `json:"system_events,omitempty"`
	GeneralFacts    []string `json:"general_facts,omitempty"`
	Events          []string `json:"events,omitempty"`
	Traits          []string `json:"traits,omitempty"`
	Tensions        []string `json:"tensions,omitempty"`
	Callbacks       []string `json:"callbacks,omitempty"`
}

type OpenLoop struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Label       string `json:"label,omitempty"`
	Owner       string `json:"owner,omitempty"`
	NextStep    string `json:"next_step,omitempty"`
}

type MemoryCapsule struct {
	ID                  int64      `json:"id"`
	DayDate             string     `json:"day_date"`
	TimeSpan            string     `json:"time_span"`
	EpisodeID           string     `json:"episode_id"`
	Kind                string     `json:"kind"`
	SourceStartRowID    int64      `json:"source_start_row_id"`
	SourceEndRowID      int64      `json:"source_end_row_id"`
	SourceMessageCount  int        `json:"source_message_count"`
	Participants        []string   `json:"participants"`
	MainTopic           string     `json:"main_topic"`
	Mood                string     `json:"mood"`
	KeyFacts            []string   `json:"key_facts"`
	TypedFacts          TypedFacts `json:"typed_facts"`
	UnresolvedQuestions []string   `json:"unresolved_questions"`
	OpenLoops           []OpenLoop `json:"open_loops"`
	Category            string     `json:"category"`
	IsMerged            bool       `json:"is_merged"`
	CreatedAt           time.Time  `json:"created_at"`
}

func (c *MemoryCapsule) Normalize() {
	if c.Participants == nil {
		c.Participants = []string{}
	}
	if c.KeyFacts == nil {
		c.KeyFacts = []string{}
	}
	if c.UnresolvedQuestions == nil {
		c.UnresolvedQuestions = []string{}
	}
	if c.OpenLoops == nil {
		c.OpenLoops = []OpenLoop{}
	}
}

func (c *MemoryCapsule) SearchBlob() string {
	return c.MainTopic + " " + strings.Join(c.KeyFacts, " ")
}

// SaveMemoryCapsule inserts a new memory capsule.
func (r *Repository) SaveMemoryCapsule(capsule *MemoryCapsule) error {
	return r.saveMemoryCapsule(r.db, capsule)
}

// SaveMemoryCapsuleTx inserts a new memory capsule inside an existing transaction.
func (r *Repository) SaveMemoryCapsuleTx(tx *sql.Tx, capsule *MemoryCapsule) error {
	return r.saveMemoryCapsule(tx, capsule)
}

func (r *Repository) saveMemoryCapsule(exec executor, capsule *MemoryCapsule) error {
	capsule.Normalize()

	participantsJSON, _ := json.Marshal(capsule.Participants)
	keyFactsJSON, _ := json.Marshal(capsule.KeyFacts)
	typedFactsJSON, _ := json.Marshal(capsule.TypedFacts)
	unresolvedQuestionsJSON, _ := json.Marshal(capsule.UnresolvedQuestions)
	openLoopsJSON, _ := json.Marshal(capsule.OpenLoops)

	res, err := exec.Exec(`
		INSERT INTO memory_capsules (
			day_date, time_span, episode_id, kind, source_start_row_id, source_end_row_id, source_message_count,
			participants, main_topic, mood, key_facts, typed_facts_json, unresolved_questions, open_loops_json,
			category, is_merged, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		capsule.DayDate,
		capsule.TimeSpan,
		capsule.EpisodeID,
		capsule.Kind,
		capsule.SourceStartRowID,
		capsule.SourceEndRowID,
		capsule.SourceMessageCount,
		participantsJSON,
		capsule.MainTopic,
		capsule.Mood,
		keyFactsJSON,
		typedFactsJSON,
		unresolvedQuestionsJSON,
		openLoopsJSON,
		capsule.Category,
		capsule.IsMerged,
		timeutil.Now(),
	)
	if err != nil {
		return err
	}
	capsule.ID, _ = res.LastInsertId()
	return nil
}

// GetMemoryCapsulesByDate retrieves all capsules for a specific date (YYYY-MM-DD).
func (r *Repository) GetMemoryCapsulesByDate(dayDate string) ([]MemoryCapsule, error) {
	rows, err := r.db.Query(`
		SELECT id, day_date, time_span, episode_id, kind, source_start_row_id, source_end_row_id, source_message_count,
		       participants, main_topic, mood, key_facts, typed_facts_json, unresolved_questions, open_loops_json,
		       category, is_merged, created_at
		FROM memory_capsules
		WHERE day_date = ?
		ORDER BY created_at ASC
	`, dayDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var capsules []MemoryCapsule
	for rows.Next() {
		capsule, err := scanMemoryCapsule(rows.Scan)
		if err != nil {
			return nil, err
		}
		capsules = append(capsules, capsule)
	}
	return capsules, nil
}

// GetActiveMemoryCapsulesByDate retrieves all unmerged capsules for a specific date (YYYY-MM-DD).
func (r *Repository) GetActiveMemoryCapsulesByDate(dayDate string) ([]MemoryCapsule, error) {
	rows, err := r.db.Query(`
		SELECT id, day_date, time_span, episode_id, kind, source_start_row_id, source_end_row_id, source_message_count,
		       participants, main_topic, mood, key_facts, typed_facts_json, unresolved_questions, open_loops_json,
		       category, is_merged, created_at
		FROM memory_capsules
		WHERE day_date = ? AND is_merged = 0
		ORDER BY created_at ASC
	`, dayDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var capsules []MemoryCapsule
	for rows.Next() {
		capsule, err := scanMemoryCapsule(rows.Scan)
		if err != nil {
			return nil, err
		}
		capsules = append(capsules, capsule)
	}
	return capsules, nil
}

// MarkCapsulesAsMerged marks a set of capsules as merged.
func (r *Repository) MarkCapsulesAsMerged(ids []int64) error {
	return r.markCapsulesAsMerged(r.db, ids)
}

// MarkCapsulesAsMergedTx marks a set of capsules as merged inside an existing transaction.
func (r *Repository) MarkCapsulesAsMergedTx(tx *sql.Tx, ids []int64) error {
	return r.markCapsulesAsMerged(tx, ids)
}

func (r *Repository) markCapsulesAsMerged(exec executor, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	for _, id := range ids {
		if _, err := exec.Exec("UPDATE memory_capsules SET is_merged = 1 WHERE id = ?", id); err != nil {
			return err
		}
	}
	return nil
}

// GetUnmergedCapsulesCount returns the number of unmerged capsules for a date.
func (r *Repository) GetUnmergedCapsulesCount(dayDate string) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM memory_capsules WHERE day_date = ? AND is_merged = 0", dayDate).Scan(&count)
	return count, err
}

// SearchPastCapsules returns at most limit relevant capsules before the provided day.
func (r *Repository) SearchPastCapsules(keywords []string, beforeDay string, limit int) ([]MemoryCapsule, error) {
	if len(keywords) == 0 || limit <= 0 {
		return nil, nil
	}

	query := ""
	for i, kw := range keywords {
		if i > 0 {
			query += " OR "
		}
		clean := strings.ReplaceAll(kw, "\"", "")
		query += "\"" + clean + "\""
	}

	rows, err := r.db.Query(`
		SELECT m.id, m.day_date, m.time_span, m.episode_id, m.kind, m.source_start_row_id, m.source_end_row_id, m.source_message_count,
		       m.participants, m.main_topic, m.mood, m.key_facts, m.typed_facts_json, m.unresolved_questions, m.open_loops_json,
		       m.category, m.is_merged, m.created_at
		FROM memory_search s
		JOIN memory_capsules m ON m.id = s.rowid
		WHERE memory_search MATCH ? AND m.day_date < ?
		ORDER BY bm25(memory_search) ASC, m.created_at DESC
		LIMIT ?
	`, query, beforeDay, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var capsules []MemoryCapsule
	for rows.Next() {
		capsule, err := scanMemoryCapsule(rows.Scan)
		if err != nil {
			return nil, err
		}
		capsules = append(capsules, capsule)
	}
	return capsules, nil
}

// GetCapsuleHWM retrieves the last encapsulated message ID for a channel.
func (r *Repository) GetCapsuleHWM(channelID string) (int64, error) {
	var lastID int64
	err := r.db.QueryRow("SELECT last_message_id FROM capsule_hwm WHERE channel_id = ?", channelID).Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return lastID, nil
}

// UpdateCapsuleHWM updates the last encapsulated message ID for a channel.
func (r *Repository) UpdateCapsuleHWM(channelID string, lastID int64) error {
	return r.updateCapsuleHWM(r.db, channelID, lastID)
}

// UpdateCapsuleHWMTx updates the last encapsulated message ID for a channel inside a transaction.
func (r *Repository) UpdateCapsuleHWMTx(tx *sql.Tx, channelID string, lastID int64) error {
	return r.updateCapsuleHWM(tx, channelID, lastID)
}

func (r *Repository) updateCapsuleHWM(exec executor, channelID string, lastID int64) error {
	_, err := exec.Exec(`
		INSERT INTO capsule_hwm (channel_id, last_message_id, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(channel_id) DO UPDATE SET
			last_message_id = excluded.last_message_id,
			updated_at = excluded.updated_at
	`, channelID, lastID)
	return err
}

// GetCapsuleHWMTx retrieves the last encapsulated message ID for a channel inside a transaction.
func (r *Repository) GetCapsuleHWMTx(tx *sql.Tx, channelID string) (int64, error) {
	var lastID int64
	err := tx.QueryRow("SELECT last_message_id FROM capsule_hwm WHERE channel_id = ?", channelID).Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return lastID, nil
}

func scanMemoryCapsule(scanFn func(dest ...any) error) (MemoryCapsule, error) {
	var capsule MemoryCapsule
	var participantsJSON, keyFactsJSON, typedFactsJSON, unresolvedQuestionsJSON, openLoopsJSON string

	err := scanFn(
		&capsule.ID,
		&capsule.DayDate,
		&capsule.TimeSpan,
		&capsule.EpisodeID,
		&capsule.Kind,
		&capsule.SourceStartRowID,
		&capsule.SourceEndRowID,
		&capsule.SourceMessageCount,
		&participantsJSON,
		&capsule.MainTopic,
		&capsule.Mood,
		&keyFactsJSON,
		&typedFactsJSON,
		&unresolvedQuestionsJSON,
		&openLoopsJSON,
		&capsule.Category,
		&capsule.IsMerged,
		&capsule.CreatedAt,
	)
	if err != nil {
		return MemoryCapsule{}, err
	}

	_ = json.Unmarshal([]byte(participantsJSON), &capsule.Participants)
	_ = json.Unmarshal([]byte(keyFactsJSON), &capsule.KeyFacts)
	_ = json.Unmarshal([]byte(typedFactsJSON), &capsule.TypedFacts)
	_ = json.Unmarshal([]byte(unresolvedQuestionsJSON), &capsule.UnresolvedQuestions)
	_ = json.Unmarshal([]byte(openLoopsJSON), &capsule.OpenLoops)
	capsule.Normalize()

	return capsule, nil
}
