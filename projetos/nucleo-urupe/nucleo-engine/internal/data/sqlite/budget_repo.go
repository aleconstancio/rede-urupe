/*
 * Copyright (c) 2026 Talos V2 Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"time"
	"nucleo-engine/internal/pkg/timeutil"
)

type BudgetEvent struct {
	ID           int64
	TurnID       string
	CostTokens   int
	InputTokens  int
	OutputTokens int
	Model        string
	Reason       string
	TriggerType  string
	CreatedAt    time.Time
}

// LogBudgetEvent records a cost event in the ledger.
func (r *Repository) LogBudgetEvent(turnID, reason, model, triggerType string, inputTokens, outputTokens int) error {
	costTokens := inputTokens + outputTokens
	_, err := r.db.Exec(`
		INSERT INTO budget_events (turn_id, cost_tokens, input_tokens, output_tokens, model, reason, trigger_type, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, turnID, costTokens, inputTokens, outputTokens, model, reason, triggerType, timeutil.Now())
	return err
}

// GetTotalCostForToday returns the sum of cost_tokens for the current Brasilia day.
func (r *Repository) GetTotalCostForToday() (int, error) {
	today := timeutil.Today()
	var total int
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(cost_tokens), 0)
		FROM budget_events
		WHERE strftime('%Y-%m-%d', created_at) = ?
	`, today).Scan(&total)
	return total, err
}

// GetTotalCostLast24Hours returns the sum of cost_tokens for the last 24 hours.
func (r *Repository) GetTotalCostLast24Hours() (int, error) {
	var total int
	since := timeutil.Now().Add(-24 * time.Hour)
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(cost_tokens), 0)
		FROM budget_events
		WHERE created_at > ?
	`, since).Scan(&total)
	return total, err
}
