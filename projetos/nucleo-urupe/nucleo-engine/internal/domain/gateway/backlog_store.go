package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/aleconstancio/talos/v2/engine/core/agent"
	"nucleo-engine/internal/data/sqlite"
)

// SQLiteBacklogStore implements agent.BacklogStore backed by SQLite.
type SQLiteBacklogStore struct {
	Repo *sqlite.Repository
}

func (s *SQLiteBacklogStore) SavePendingTurn(ctx context.Context, turn *agent.PendingTurn) error {
	_, err := s.Repo.GetDB().Exec(`
		INSERT INTO pending_turns (channel_id, message_id, content, retry_count, status, error, created_at, updated_at)
		VALUES (?, ?, ?, 0, 'pending', '', ?, ?)
	`, turn.ChannelID, turn.MessageID, turn.Content,
		turn.CreatedAt.Format("2006-01-02 15:04:05"),
		turn.UpdatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return fmt.Errorf("save pending turn: %w", err)
	}
	return nil
}

func (s *SQLiteBacklogStore) GetPendingTurns(ctx context.Context, limit int) ([]*agent.PendingTurn, error) {
	rows, err := s.Repo.GetDB().Query(`
		SELECT id, channel_id, message_id, content, retry_count, status, error, created_at, updated_at
		FROM pending_turns WHERE status = 'pending' ORDER BY created_at ASC LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turns []*agent.PendingTurn
	for rows.Next() {
		t := &agent.PendingTurn{}
		var createdAt, updatedAt string
		err := rows.Scan(&t.ID, &t.ChannelID, &t.MessageID, &t.Content, &t.RetryCount, &t.Status, &t.Error, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
		t.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		t.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		turns = append(turns, t)
	}
	return turns, nil
}

func (s *SQLiteBacklogStore) MarkCompleted(ctx context.Context, turnID int64) error {
	_, err := s.Repo.GetDB().Exec("UPDATE pending_turns SET status = 'completed', updated_at = ? WHERE id = ?",
		time.Now().UTC().Format("2006-01-02 15:04:05"), turnID)
	return err
}

func (s *SQLiteBacklogStore) MarkFailed(ctx context.Context, turnID int64, errMsg string) error {
	_, err := s.Repo.GetDB().Exec("UPDATE pending_turns SET status = 'failed', error = ?, retry_count = retry_count + 1, updated_at = ? WHERE id = ?",
		errMsg, time.Now().UTC().Format("2006-01-02 15:04:05"), turnID)
	return err
}
