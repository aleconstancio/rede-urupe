/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

type SentimentDaily struct {
	Date          string  `json:"date"`
	ChannelID     string  `json:"channel_id"`
	AvgIntensity  float64 `json:"avg_intensity"`
	AvgEuphoria   float64 `json:"avg_euphoria"`
	AvgConflict   float64 `json:"avg_conflict"`
	MessageCount  int     `json:"message_count"`
}

type ChannelHealth struct {
	ChannelID      string  `json:"channel_id"`
	ChannelName    string  `json:"channel_name"`
	TotalMessages  int     `json:"total_messages"`
	UniqueAuthors  int     `json:"unique_authors"`
	BotMessages    int     `json:"bot_messages"`
	AvgReactions   float64 `json:"avg_reactions"`
	MessageVelocity float64 `json:"message_velocity"` // msgs per hour
}

type TokenUsageByModel struct {
	Model    string `json:"model"`
	TotalTokens int  `json:"total_tokens"`
	CallCount int    `json:"call_count"`
}

type TokenUsageByReason struct {
	Reason    string `json:"reason"`
	TotalTokens int  `json:"total_tokens"`
	CallCount int    `json:"call_count"`
}

type TokenUsageTrend struct {
	Date        string `json:"date"`
	TotalTokens int    `json:"total_tokens"`
	CallCount   int    `json:"call_count"`
}

func (r *Repository) GetSentimentTrend(channelID string, days int) ([]SentimentDaily, error) {
	rows, err := r.db.Query(`
		SELECT DATE(timestamp) as day,
			CASE
				WHEN AVG(LENGTH(content)) > 100 THEN 0.3
				WHEN AVG(LENGTH(content)) > 50 THEN 0.2
				ELSE 0.1
			END as avg_intensity,
			0.0 as avg_euphoria,
			CASE
				WHEN SUM(CASE WHEN content LIKE '%porra%' OR content LIKE '%caralho%' OR content LIKE '%merda%' THEN 1 ELSE 0 END) > 0 THEN 0.3
				ELSE 0.0
			END as avg_conflict,
			COUNT(*) as msg_count
		FROM messages
		WHERE channel_id = ? AND is_bot = 0 AND is_deleted = 0
			AND DATE(timestamp) >= DATE('now', '-' || ? || ' days')
		GROUP BY DATE(timestamp)
		ORDER BY day ASC
	`, channelID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SentimentDaily
	for rows.Next() {
		var s SentimentDaily
		if err := rows.Scan(&s.Date, &s.AvgIntensity, &s.AvgEuphoria, &s.AvgConflict, &s.MessageCount); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (r *Repository) GetChannelHealth(channelID string) (*ChannelHealth, error) {
	h := &ChannelHealth{ChannelID: channelID}

	err := r.db.QueryRow(`
		SELECT COUNT(*), COUNT(DISTINCT author_id),
			SUM(CASE WHEN is_bot THEN 1 ELSE 0 END)
		FROM messages WHERE channel_id = ? AND is_deleted = 0
	`, channelID).Scan(&h.TotalMessages, &h.UniqueAuthors, &h.BotMessages)
	if err != nil {
		return nil, err
	}

	_ = r.db.QueryRow(`
		SELECT AVG(reaction_count) FROM (
			SELECT COUNT(*) as reaction_count
			FROM messages, json_each(reactions_json)
			WHERE channel_id = ? AND is_deleted = 0
			GROUP BY id
		)
	`, channelID).Scan(&h.AvgReactions)

	return h, nil
}

func (r *Repository) GetTokenUsageByModel(days int) ([]TokenUsageByModel, error) {
	rows, err := r.db.Query(`
		SELECT model, SUM(cost_tokens) as total_tokens, COUNT(*) as call_count
		FROM budget_events
		WHERE created_at >= datetime('now', '-' || ? || ' days')
		GROUP BY model
		ORDER BY total_tokens DESC
	`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TokenUsageByModel
	for rows.Next() {
		var t TokenUsageByModel
		if err := rows.Scan(&t.Model, &t.TotalTokens, &t.CallCount); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *Repository) GetTokenUsageByReason(days int) ([]TokenUsageByReason, error) {
	rows, err := r.db.Query(`
		SELECT reason, SUM(cost_tokens) as total_tokens, COUNT(*) as call_count
		FROM budget_events
		WHERE created_at >= datetime('now', '-' || ? || ' days')
		GROUP BY reason
		ORDER BY total_tokens DESC
	`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TokenUsageByReason
	for rows.Next() {
		var t TokenUsageByReason
		if err := rows.Scan(&t.Reason, &t.TotalTokens, &t.CallCount); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *Repository) GetTokenUsageTrend(days int) ([]TokenUsageTrend, error) {
	rows, err := r.db.Query(`
		SELECT DATE(created_at) as day, SUM(cost_tokens) as total_tokens, COUNT(*) as call_count
		FROM budget_events
		WHERE created_at >= datetime('now', '-' || ? || ' days')
		GROUP BY DATE(created_at)
		ORDER BY day ASC
	`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TokenUsageTrend
	for rows.Next() {
		var t TokenUsageTrend
		if err := rows.Scan(&t.Date, &t.TotalTokens, &t.CallCount); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (r *Repository) GetEngagementScore(channelID string, days int) (float64, error) {
	var score float64
	err := r.db.QueryRow(`
		SELECT CAST(COUNT(DISTINCT author_id) AS REAL) * 10.0 /
			MAX(CAST((julianday('now') - julianday(MIN(timestamp))) AS REAL) * 24.0, 1.0)
		FROM messages
		WHERE channel_id = ? AND is_bot = 0 AND is_deleted = 0
			AND timestamp >= datetime('now', '-' || ? || ' days')
	`, channelID, days).Scan(&score)
	return score, err
}
