/*
 * Copyright (c) 2026 Talos V2 Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/aleconstancio/talos/v2/behavior/identity"
	"nucleo-engine/internal/pkg/timeutil"
)

// SaveStanceEvent persists a new stance event to the ledger.
func (r *Repository) SaveStanceEvent(tx *sql.Tx, channelID, authorID, topic, position, action, evidenceType string, confidence float64, sourceMsgID string) error {
	exec := r.getExecutor(tx)
	channelID = strings.TrimSpace(channelID)
	authorID = strings.TrimSpace(authorID)
	topic = strings.TrimSpace(topic)
	position = strings.TrimSpace(position)
	action = strings.TrimSpace(action)
	evidenceType = strings.TrimSpace(evidenceType)
	sourceMsgID = strings.TrimSpace(sourceMsgID)

	if channelID == "" || authorID == "" || topic == "" || position == "" || action == "" || evidenceType == "" {
		return nil
	}

	var existingID int64
	err := exec.QueryRow(`
		SELECT id
		FROM stance_events
		WHERE channel_id = ?
		  AND author_id = ?
		  AND topic = ?
		  AND position = ?
		  AND action = ?
		  AND evidence_type = ?
		  AND (
		    (? <> '' AND source_msg_id = ?)
		    OR
		    (? = '' AND source_msg_id = '' AND created_at >= ?)
		  )
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, channelID, authorID, topic, position, action, evidenceType, sourceMsgID, sourceMsgID, sourceMsgID, timeutil.Now().Add(-6*time.Hour)).Scan(&existingID)
	if err == nil {
		_, err = exec.Exec(`
			UPDATE stance_events
			SET confidence = CASE WHEN confidence < ? THEN ? ELSE confidence END,
			    validated = 1,
			    source_msg_id = CASE WHEN source_msg_id = '' THEN ? ELSE source_msg_id END
			WHERE id = ?
		`, confidence, confidence, sourceMsgID, existingID)
		return err
	}
	if err != sql.ErrNoRows {
		return err
	}

	_, err = exec.Exec(`
		INSERT INTO stance_events (channel_id, author_id, topic, position, action, evidence_type, confidence, source_msg_id, episode_id, validated, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 0, 1, ?)
	`, channelID, authorID, topic, position, action, evidenceType, confidence, sourceMsgID, timeutil.Now())
	return err
}

// SaveBehaviorEvent persists a new behavioral observation.
func (r *Repository) SaveBehaviorEvent(tx *sql.Tx, channelID, authorID, archetype, style string) error {
	exec := r.getExecutor(tx)
	channelID = strings.TrimSpace(channelID)
	authorID = strings.TrimSpace(authorID)
	archetype = strings.TrimSpace(archetype)
	style = strings.TrimSpace(style)
	if channelID == "" || authorID == "" || archetype == "" || style == "" {
		return nil
	}

	var existingID int64
	err := exec.QueryRow(`
		SELECT id
		FROM behavior_events
		WHERE channel_id = ?
		  AND author_id = ?
		  AND archetype = ?
		  AND rhetorical_style = ?
		  AND created_at >= ?
		ORDER BY created_at DESC, id DESC
		LIMIT 1
	`, channelID, authorID, archetype, style, timeutil.Now().Add(-24*time.Hour)).Scan(&existingID)
	if err == nil {
		_, err = exec.Exec(`UPDATE behavior_events SET status = 'approved' WHERE id = ?`, existingID)
		return err
	}
	if err != sql.ErrNoRows {
		return err
	}

	_, err = exec.Exec(`
		INSERT INTO behavior_events (channel_id, author_id, archetype, rhetorical_style, status, created_at)
		VALUES (?, ?, ?, ?, 'approved', ?)
	`, channelID, authorID, archetype, style, timeutil.Now())
	return err
}

// GetParticipantProjections resolves stable participant summaries from stance and behavior events.
func (r *Repository) GetParticipantProjections(channelID string, authorIDs []string, now time.Time, limit int) ([]identity.ParticipantProjection, error) {
	ids := uniqueAuthorIDs(authorIDs, limit)
	if len(ids) == 0 {
		return nil, nil
	}

	displayNames, err := r.getLatestDisplayNames(channelID, ids)
	if err != nil {
		return nil, err
	}
	stanceEvents, err := r.GetStanceEvents(ids, channelID)
	if err != nil {
		return nil, err
	}
	behaviorEvents, err := r.GetBehaviorEvents(ids, channelID)
	if err != nil {
		return nil, err
	}

	stanceMap := make(map[string][]identity.StanceEvent, len(ids))
	for _, event := range stanceEvents {
		stanceMap[event.ParticipantID] = append(stanceMap[event.ParticipantID], event)
	}

	behaviorMap := make(map[string][]identity.BehaviorEvent, len(ids))
	for _, event := range behaviorEvents {
		behaviorMap[event.ParticipantID] = append(behaviorMap[event.ParticipantID], event)
	}

	projections := make([]identity.ParticipantProjection, 0, len(ids))
	for _, id := range ids {
		displayName := strings.TrimSpace(displayNames[id])
		if displayName == "" {
			displayName = id
		}
		projections = append(projections, identity.Project(
			id,
			displayName,
			stanceMap[id],
			behaviorMap[id],
			now,
		))
	}
	return projections, nil
}

func (r *Repository) GetStanceEvents(authorIDs []string, channelID string) ([]identity.StanceEvent, error) {
	if len(authorIDs) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(authorIDs))
	args := make([]interface{}, len(authorIDs)+1)
	args[0] = channelID
	for i, id := range authorIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := fmt.Sprintf(`
		SELECT id, channel_id, author_id, topic, position, action, evidence_type, confidence, source_msg_id, episode_id, validated, created_at
		FROM stance_events
		WHERE channel_id = ? AND author_id IN (%s)
		ORDER BY created_at DESC`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []identity.StanceEvent
	for rows.Next() {
		var e identity.StanceEvent
		if err := rows.Scan(&e.ID, &e.ChannelID, &e.ParticipantID, &e.Topic, &e.Position, &e.Action, &e.EvidenceType, &e.Confidence, &e.SourceMsgID, &e.EpisodeID, &e.Validated, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *Repository) GetBehaviorEvents(authorIDs []string, channelID string) ([]identity.BehaviorEvent, error) {
	if len(authorIDs) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(authorIDs))
	args := make([]interface{}, len(authorIDs)+1)
	args[0] = channelID
	for i, id := range authorIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := fmt.Sprintf(`
		SELECT id, channel_id, author_id, episode_id, archetype, rhetorical_style, status, created_at
		FROM behavior_events
		WHERE channel_id = ? AND author_id IN (%s)
		ORDER BY created_at DESC`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []identity.BehaviorEvent
	for rows.Next() {
		var e identity.BehaviorEvent
		if err := rows.Scan(&e.ID, &e.ChannelID, &e.ParticipantID, &e.EpisodeID, &e.Archetype, &e.RhetoricalStyle, &e.Status, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}
func (r *Repository) getLatestDisplayNames(channelID string, authorIDs []string) (map[string]string, error) {
	if len(authorIDs) == 0 {
		return make(map[string]string), nil
	}
	placeholders := make([]string, len(authorIDs))
	args := make([]interface{}, len(authorIDs)+1)
	args[0] = channelID
	for i, id := range authorIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}

	query := fmt.Sprintf(`
		SELECT author_id, author
		FROM (
			SELECT author_id, author, ROW_NUMBER() OVER (PARTITION BY author_id ORDER BY timestamp DESC) as rn
			FROM messages
			WHERE channel_id = ? AND author_id IN (%s)
		)
		WHERE rn = 1
	`, strings.Join(placeholders, ","))

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]string)
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		out[id] = name
	}
	return out, nil
}

