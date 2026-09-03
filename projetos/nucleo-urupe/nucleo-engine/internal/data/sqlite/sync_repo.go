/*
 * Copyright (c) 2026 Talos V2 Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"database/sql"
	"nucleo-engine/internal/pkg/timeutil"
)


// GetSyncState retrieves the current synchronization state for a channel.
func (r *Repository) GetSyncState(channelID string) (latestID, oldestID string, err error) {
	err = r.db.QueryRow("SELECT latest_synced_msg_id, oldest_synced_msg_id FROM sync_state WHERE channel_id = ?", channelID).Scan(&latestID, &oldestID)
	if err == sql.ErrNoRows {
		return "", "", nil
	}
	return
}

// UpdateSyncState updates the synchronization state for a channel.
func (r *Repository) UpdateSyncState(channelID, latestID, oldestID string) error {
	_, err := r.db.Exec(`
		INSERT INTO sync_state (channel_id, latest_synced_msg_id, oldest_synced_msg_id, updated_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(channel_id) DO UPDATE SET
			latest_synced_msg_id = excluded.latest_synced_msg_id,
			oldest_synced_msg_id = excluded.oldest_synced_msg_id,
			updated_at = excluded.updated_at
	`, channelID, latestID, oldestID, timeutil.Now())
	return err
}
