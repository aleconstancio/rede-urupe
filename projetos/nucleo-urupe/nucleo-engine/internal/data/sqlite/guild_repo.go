package sqlite

import "nucleo-engine/internal/pkg/timeutil"

type GuildConfig struct {
	GuildID          string `json:"guild_id"`
	Prefix           string `json:"prefix"`
	DefaultChannelID string `json:"default_channel_id"`
	ModLogChannelID  string `json:"mod_log_channel_id"`
	GateModel        string `json:"gate_model"`
	ReplyModel       string `json:"reply_model"`
	MemoryModel      string `json:"memory_model"`
	ModEnabled       bool   `json:"mod_enabled"`
	MaxWarnings      int    `json:"max_warnings"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

func (r *Repository) GetGuildConfig(guildID string) (*GuildConfig, error) {
	gc := &GuildConfig{}
	err := r.db.QueryRow(`
		SELECT guild_id, prefix, default_channel_id, mod_log_channel_id,
		       gate_model, reply_model, memory_model, mod_enabled, max_warnings,
		       created_at, updated_at
		FROM guild_config WHERE guild_id = ?
	`, guildID).Scan(
		&gc.GuildID, &gc.Prefix, &gc.DefaultChannelID, &gc.ModLogChannelID,
		&gc.GateModel, &gc.ReplyModel, &gc.MemoryModel, &gc.ModEnabled, &gc.MaxWarnings,
		&gc.CreatedAt, &gc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return gc, nil
}

func (r *Repository) UpsertGuildConfig(gc GuildConfig) error {
	now := timeutil.Now()
	_, err := r.db.Exec(`
		INSERT INTO guild_config (guild_id, prefix, default_channel_id, mod_log_channel_id,
		                          gate_model, reply_model, memory_model, mod_enabled, max_warnings,
		                          created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(guild_id) DO UPDATE SET
			prefix = excluded.prefix,
			default_channel_id = excluded.default_channel_id,
			mod_log_channel_id = excluded.mod_log_channel_id,
			gate_model = excluded.gate_model,
			reply_model = excluded.reply_model,
			memory_model = excluded.memory_model,
			mod_enabled = excluded.mod_enabled,
			max_warnings = excluded.max_warnings,
			updated_at = excluded.updated_at
	`, gc.GuildID, gc.Prefix, gc.DefaultChannelID, gc.ModLogChannelID,
		gc.GateModel, gc.ReplyModel, gc.MemoryModel, gc.ModEnabled, gc.MaxWarnings,
		now, now)
	return err
}

func (r *Repository) ListGuilds() ([]string, error) {
	rows, err := r.db.Query("SELECT guild_id FROM guild_config")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repository) AddWarning(guildID, userID, channelID, reason, messageID string, severity int) error {
	_, err := r.db.Exec(`
		INSERT INTO guild_warnings (guild_id, user_id, channel_id, reason, message_id, severity)
		VALUES (?, ?, ?, ?, ?, ?)
	`, guildID, userID, channelID, reason, messageID, severity)
	return err
}

func (r *Repository) GetWarningCount(guildID, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM guild_warnings WHERE guild_id = ? AND user_id = ?",
		guildID, userID,
	).Scan(&count)
	return count, err
}

func (r *Repository) ListActiveChannels(guildID string) ([]string, error) {
	rows, err := r.db.Query("SELECT channel_id FROM active_channels WHERE guild_id = ? AND is_active = 1", guildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var channels []string
	for rows.Next() {
		var ch string
		if err := rows.Scan(&ch); err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}
	return channels, nil
}

func (r *Repository) RegisterChannel(guildID, channelID string) error {
	_, err := r.db.Exec(
		"INSERT OR IGNORE INTO active_channels (guild_id, channel_id) VALUES (?, ?)",
		guildID, channelID,
	)
	return err
}
