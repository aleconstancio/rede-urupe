package sqlite

import (
	"time"
)

type PendingTurn struct {
	ID           int64     `json:"id"`
	ChannelID    string    `json:"channel_id"`
	MessageID    string    `json:"message_id"`
	MessageRowID int64     `json:"message_row_id"`
	IsReactive   bool      `json:"is_reactive"`
	Reason       string    `json:"reason"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	Timestamp    time.Time `json:"timestamp"`
}

func (r *Repository) SavePendingTurn(channelID, messageID string, messageRowID int64, timestamp time.Time, isReactive bool, reason string) error {
	_, err := r.db.Exec(`
		INSERT INTO pending_turns (channel_id, message_id, message_row_id, is_reactive, reason, status, created_at)
		VALUES (?, ?, ?, ?, ?, 'pending', ?)
	`, channelID, messageID, messageRowID, isReactive, reason, timestamp)
	return err
}

func (r *Repository) GetPendingTurns(channelID string, limit ...int) ([]PendingTurn, error) {
	query := `
		SELECT id, channel_id, message_id, message_row_id, is_reactive, reason, status, created_at
		FROM pending_turns
		WHERE channel_id = ? AND status = 'pending'
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turns []PendingTurn
	for rows.Next() {
		var t PendingTurn
		if err := rows.Scan(&t.ID, &t.ChannelID, &t.MessageID, &t.MessageRowID, &t.IsReactive, &t.Reason, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.Timestamp = t.CreatedAt
		turns = append(turns, t)
	}
	return turns, nil
}

func (r *Repository) DeletePendingTurn(id int64) error {
	_, err := r.db.Exec("DELETE FROM pending_turns WHERE id = ?", id)
	return err
}

func (r *Repository) UpdatePendingTurnStatus(id int64, status string, errMsg ...string) error {
	_, err := r.db.Exec("UPDATE pending_turns SET status = ? WHERE id = ?", status, id)
	return err
}
