package sqlite

import (
	"time"
)

type GatewayState struct {
	ID                 int64     `json:"id"`
	ChannelID          string    `json:"channel_id"`
	LastPulseRowID     int64     `json:"last_pulse_row_id"`
	LastPulseMsgID     string    `json:"last_pulse_msg_id"`
	LastPulseAt        time.Time `json:"last_pulse_at"`
	LastPulseTimestamp time.Time `json:"last_pulse_timestamp"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (r *Repository) GetGatewayState(channelID string) (*GatewayState, error) {
	var state GatewayState
	err := r.db.QueryRow(`
		SELECT id, channel_id, last_pulse_row_id, COALESCE(last_pulse_msg_id, ''), updated_at
		FROM gateway_state
		WHERE channel_id = ?
	`, channelID).Scan(&state.ID, &state.ChannelID, &state.LastPulseRowID, &state.LastPulseMsgID, &state.UpdatedAt)
	if err != nil {
		return &GatewayState{ChannelID: channelID, LastPulseRowID: 0, UpdatedAt: time.Now(), LastPulseAt: time.Now(), LastPulseTimestamp: time.Now()}, nil
	}
	state.LastPulseAt = state.UpdatedAt
	state.LastPulseTimestamp = state.UpdatedAt
	return &state, nil
}

func (r *Repository) SaveGatewayState(state *GatewayState) error {
	_, err := r.db.Exec(`
		INSERT INTO gateway_state (channel_id, last_pulse_row_id, last_pulse_msg_id, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(channel_id) DO UPDATE SET
			last_pulse_row_id = excluded.last_pulse_row_id,
			last_pulse_msg_id = excluded.last_pulse_msg_id,
			updated_at = excluded.updated_at
	`, state.ChannelID, state.LastPulseRowID, state.LastPulseMsgID)
	return err
}

func (r *Repository) UpsertGatewayState(state *GatewayState) error {
	return r.SaveGatewayState(state)
}

func (r *Repository) CountMessagesSince(channelID string, since time.Time) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM messages WHERE channel_id = ? AND timestamp > ?", channelID, since).Scan(&count)
	return count, err
}
