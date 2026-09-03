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
)

// Message represents a single Discord message log entry.
type Message struct {
	ID                int64     `json:"id"`
	DiscordID         string    `json:"discord_id"`
	Timestamp         time.Time `json:"timestamp"`
	Author            string    `json:"author"`
	AuthorID          string    `json:"author_id"`
	Content           string    `json:"content"`
	ChannelID         string    `json:"channel_id"`
	Category          string    `json:"category"`
	Mentions          []string  `json:"mentions"`
	Attachments       []string  `json:"attachments"`
	Reactions         []string  `json:"reactions"`
	ReplyToID         string    `json:"reply_to_id"`
	EpisodeID         string    `json:"episode_id"`
	IsBot             bool      `json:"is_bot"`
	IsDeleted         bool      `json:"is_deleted"`
	InternalMonologue string    `json:"internal_monologue"`
	GroundingLedger   string    `json:"grounding_ledger"`
}

// MarkMessageDeleted marks a message as deleted in the log.
func (r *Repository) MarkMessageDeleted(channelID, messageID string, deletedAt time.Time) error {
	_, err := r.db.Exec("UPDATE messages SET is_deleted = 1 WHERE channel_id = ? AND discord_id = ?", channelID, messageID)
	return err
}

// GetRecentMessages fetches the last N messages for a channel.
func (r *Repository) GetRecentMessages(channelID string, limit int) ([]Message, error) {
	rows, err := r.db.Query(`
		SELECT id, discord_id, author, author_id, content, channel_id, category, mentions_json, attachments_json, reactions_json, is_bot, timestamp, reply_to_id, episode_id, internal_monologue, grounding_ledger
		FROM messages
		WHERE channel_id = ? AND is_deleted = 0
		ORDER BY timestamp DESC
		LIMIT ?
	`, channelID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		var replyTo sql.NullString
		var mentionsJSON, attachmentsJSON, reactionsJSON string
		err := rows.Scan(&m.ID, &m.DiscordID, &m.Author, &m.AuthorID, &m.Content, &m.ChannelID, &m.Category, &mentionsJSON, &attachmentsJSON, &reactionsJSON, &m.IsBot, &m.Timestamp, &replyTo, &m.EpisodeID, &m.InternalMonologue, &m.GroundingLedger)
		if err != nil {
			return nil, err
		}
		m.Mentions = decodeStringSlice(mentionsJSON)
		m.Attachments = decodeStringSlice(attachmentsJSON)
		m.Reactions = decodeStringSlice(reactionsJSON)
		m.ReplyToID = replyTo.String
		msgs = append(msgs, m)
	}
	// Keep newest at the top for the Command Center feel
	return msgs, nil
}

// GetMessagesByTimeRange fetches messages for a channel within a specific time range.
func (r *Repository) GetMessagesByTimeRange(channelID string, start, end time.Time) ([]Message, error) {
	rows, err := r.db.Query(`
		SELECT id, discord_id, author, author_id, content, channel_id, category, mentions_json, attachments_json, reactions_json, is_bot, timestamp, reply_to_id, episode_id, internal_monologue, grounding_ledger
		FROM messages
		WHERE channel_id = ? AND timestamp BETWEEN ? AND ? AND is_deleted = 0
		ORDER BY timestamp ASC
	`, channelID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		var replyTo sql.NullString
		var mentionsJSON, attachmentsJSON, reactionsJSON string
		err := rows.Scan(&m.ID, &m.DiscordID, &m.Author, &m.AuthorID, &m.Content, &m.ChannelID, &m.Category, &mentionsJSON, &attachmentsJSON, &reactionsJSON, &m.IsBot, &m.Timestamp, &replyTo, &m.EpisodeID, &m.InternalMonologue, &m.GroundingLedger)
		if err != nil {
			return nil, err
		}
		m.Mentions = decodeStringSlice(mentionsJSON)
		m.Attachments = decodeStringSlice(attachmentsJSON)
		m.Reactions = decodeStringSlice(reactionsJSON)
		m.ReplyToID = replyTo.String
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// SaveMessage persists a message to the database with rich metadata and returns the local row id.
func (r *Repository) SaveMessage(discordID, author, authorID, content, channelID, category string, isBot bool, timestamp time.Time, replyToID string, mentions, attachments, reactions []string, internalMonologue, groundingLedger string) (int64, error) {
	mentionsJSON, _ := json.Marshal(mentions)
	attachmentsJSON, _ := json.Marshal(attachments)
	reactionsJSON, _ := json.Marshal(reactions)

	var id int64
	err := r.db.QueryRow(`
		INSERT INTO messages (discord_id, author, author_id, content, channel_id, category, is_bot, timestamp, reply_to_id, mentions_json, attachments_json, reactions_json, internal_monologue, grounding_ledger)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(discord_id) DO UPDATE SET
			author = excluded.author,
			author_id = excluded.author_id,
			content = excluded.content,
			channel_id = excluded.channel_id,
			category = excluded.category,
			is_bot = excluded.is_bot,
			timestamp = excluded.timestamp,
			reply_to_id = excluded.reply_to_id,
			mentions_json = excluded.mentions_json,
			attachments_json = excluded.attachments_json,
			reactions_json = excluded.reactions_json,
			internal_monologue = excluded.internal_monologue,
			grounding_ledger = excluded.grounding_ledger,
			is_deleted = 0
		RETURNING id
	`, discordID, author, authorID, content, channelID, category, isBot, timestamp, replyToID, string(mentionsJSON), string(attachmentsJSON), string(reactionsJSON), internalMonologue, groundingLedger).Scan(&id)
	return id, err
}

// GetMessagesByRowRange fetches messages for a channel within a specific local row-id window.
func (r *Repository) GetMessagesByRowRange(channelID string, startExclusive, endInclusive int64) ([]Message, error) {
	rows, err := r.db.Query(`
		SELECT id, discord_id, author, author_id, content, channel_id, category, mentions_json, attachments_json, reactions_json, is_bot, timestamp, reply_to_id, episode_id, internal_monologue, grounding_ledger
		FROM messages
		WHERE channel_id = ? AND id > ? AND id <= ? AND is_deleted = 0
		ORDER BY id ASC
	`, channelID, startExclusive, endInclusive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		var replyTo sql.NullString
		var mentionsJSON, attachmentsJSON, reactionsJSON string
		err := rows.Scan(&m.ID, &m.DiscordID, &m.Author, &m.AuthorID, &m.Content, &m.ChannelID, &m.Category, &mentionsJSON, &attachmentsJSON, &reactionsJSON, &m.IsBot, &m.Timestamp, &replyTo, &m.EpisodeID, &m.InternalMonologue, &m.GroundingLedger)
		if err != nil {
			return nil, err
		}
		m.Mentions = decodeStringSlice(mentionsJSON)
		m.Attachments = decodeStringSlice(attachmentsJSON)
		m.Reactions = decodeStringSlice(reactionsJSON)
		m.ReplyToID = replyTo.String
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// AssignEpisodeToMessagesByRowRange stamps an episode identifier onto a message span.
func (r *Repository) AssignEpisodeToMessagesByRowRange(channelID string, startExclusive, endInclusive int64, episodeID string) error {
	return r.assignEpisodeToMessagesByRowRange(r.db, channelID, startExclusive, endInclusive, episodeID)
}

// AssignEpisodeToMessagesByRowRangeTx stamps an episode identifier onto a message span inside a transaction.
func (r *Repository) AssignEpisodeToMessagesByRowRangeTx(tx *sql.Tx, channelID string, startExclusive, endInclusive int64, episodeID string) error {
	return r.assignEpisodeToMessagesByRowRange(tx, channelID, startExclusive, endInclusive, episodeID)
}

func (r *Repository) assignEpisodeToMessagesByRowRange(exec executor, channelID string, startExclusive, endInclusive int64, episodeID string) error {
	_, err := exec.Exec(`
		UPDATE messages
		SET episode_id = ?
		WHERE channel_id = ? AND id > ? AND id <= ? AND is_deleted = 0
	`, episodeID, channelID, startExclusive, endInclusive)
	return err
}

// GetLatestMessageByChannel returns the newest message for a channel.
func (r *Repository) GetLatestMessageByChannel(channelID string) (*Message, error) {
	row := r.db.QueryRow(`
		SELECT id, discord_id, author, author_id, content, channel_id, category, mentions_json, attachments_json, reactions_json, is_bot, timestamp, reply_to_id, episode_id, internal_monologue, grounding_ledger
		FROM messages
		WHERE channel_id = ? AND is_deleted = 0
		ORDER BY id DESC
		LIMIT 1
	`, channelID)

	var m Message
	var replyTo sql.NullString
	var mentionsJSON, attachmentsJSON, reactionsJSON string
	if err := row.Scan(&m.ID, &m.DiscordID, &m.Author, &m.AuthorID, &m.Content, &m.ChannelID, &m.Category, &mentionsJSON, &attachmentsJSON, &reactionsJSON, &m.IsBot, &m.Timestamp, &replyTo, &m.EpisodeID, &m.InternalMonologue, &m.GroundingLedger); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m.Mentions = decodeStringSlice(mentionsJSON)
	m.Attachments = decodeStringSlice(attachmentsJSON)
	m.Reactions = decodeStringSlice(reactionsJSON)
	m.ReplyToID = replyTo.String
	return &m, nil
}

// GetMinimumMessageRowIDByChannelSince returns the smallest local row id for non-deleted messages at or after a timestamp.
func (r *Repository) GetMinimumMessageRowIDByChannelSince(channelID string, since time.Time) (int64, error) {
	var id sql.NullInt64
	err := r.db.QueryRow(`
		SELECT MIN(id)
		FROM messages
		WHERE channel_id = ? AND is_deleted = 0 AND timestamp >= ?
	`, channelID, since).Scan(&id)
	if err != nil {
		return 0, err
	}
	if !id.Valid {
		return 0, nil
	}
	return id.Int64, nil
}

// GetLatestUnassignedMessageByChannelSince returns the newest non-deleted message without an episode assignment at or after a timestamp.
func (r *Repository) GetLatestUnassignedMessageByChannelSince(channelID string, since time.Time) (*Message, error) {
	row := r.db.QueryRow(`
		SELECT id, discord_id, author, author_id, content, channel_id, category, mentions_json, attachments_json, reactions_json, is_bot, timestamp, reply_to_id, episode_id, internal_monologue, grounding_ledger
		FROM messages
		WHERE channel_id = ? AND is_deleted = 0 AND episode_id = '' AND timestamp >= ?
		ORDER BY id DESC
		LIMIT 1
	`, channelID, since)

	var m Message
	var replyTo sql.NullString
	var mentionsJSON, attachmentsJSON, reactionsJSON string
	if err := row.Scan(&m.ID, &m.DiscordID, &m.Author, &m.AuthorID, &m.Content, &m.ChannelID, &m.Category, &mentionsJSON, &attachmentsJSON, &reactionsJSON, &m.IsBot, &m.Timestamp, &replyTo, &m.EpisodeID, &m.InternalMonologue, &m.GroundingLedger); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m.Mentions = decodeStringSlice(mentionsJSON)
	m.Attachments = decodeStringSlice(attachmentsJSON)
	m.Reactions = decodeStringSlice(reactionsJSON)
	m.ReplyToID = replyTo.String
	return &m, nil
}

// GetRecentMessagesByChannelBeforeRowSince returns up to limit messages ordered newest-first, bounded by row id and timestamp.
func (r *Repository) GetRecentMessagesByChannelBeforeRowSince(channelID string, endInclusive int64, since time.Time, limit int) ([]Message, error) {
	rows, err := r.db.Query(`
		SELECT id, discord_id, author, author_id, content, channel_id, category, mentions_json, attachments_json, reactions_json, is_bot, timestamp, reply_to_id, episode_id, internal_monologue, grounding_ledger
		FROM messages
		WHERE channel_id = ? AND is_deleted = 0 AND id <= ? AND timestamp >= ?
		ORDER BY id DESC
		LIMIT ?
	`, channelID, endInclusive, since, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		var replyTo sql.NullString
		var mentionsJSON, attachmentsJSON, reactionsJSON string
		err := rows.Scan(&m.ID, &m.DiscordID, &m.Author, &m.AuthorID, &m.Content, &m.ChannelID, &m.Category, &mentionsJSON, &attachmentsJSON, &reactionsJSON, &m.IsBot, &m.Timestamp, &replyTo, &m.EpisodeID, &m.InternalMonologue, &m.GroundingLedger)
		if err != nil {
			return nil, err
		}
		m.Mentions = decodeStringSlice(mentionsJSON)
		m.Attachments = decodeStringSlice(attachmentsJSON)
		m.Reactions = decodeStringSlice(reactionsJSON)
		m.ReplyToID = replyTo.String
		msgs = append(msgs, m)
	}
	return msgs, nil
}

// CountAssignedMessagesByRowRangeTx returns how many non-deleted messages in the row range already belong to an episode.
func (r *Repository) CountAssignedMessagesByRowRangeTx(tx *sql.Tx, channelID string, startExclusive, endInclusive int64) (int, error) {
	var count int
	err := tx.QueryRow(`
		SELECT COUNT(*)
		FROM messages
		WHERE channel_id = ? AND id > ? AND id <= ? AND is_deleted = 0 AND episode_id <> ''
	`, channelID, startExclusive, endInclusive).Scan(&count)
	return count, err
}

func (r *Repository) GetMessageByDiscordID(discordID string) (*Message, error) {
	row := r.db.QueryRow(`
		SELECT id, discord_id, author, author_id, content, channel_id, category, mentions_json, attachments_json, reactions_json, is_bot, timestamp, reply_to_id, episode_id, internal_monologue, grounding_ledger
		FROM messages
		WHERE discord_id = ?
	`, discordID)

	var m Message
	var replyTo sql.NullString
	var mentionsJSON, attachmentsJSON, reactionsJSON string
	if err := row.Scan(&m.ID, &m.DiscordID, &m.Author, &m.AuthorID, &m.Content, &m.ChannelID, &m.Category, &mentionsJSON, &attachmentsJSON, &reactionsJSON, &m.IsBot, &m.Timestamp, &replyTo, &m.EpisodeID, &m.InternalMonologue, &m.GroundingLedger); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	m.Mentions = decodeStringSlice(mentionsJSON)
	m.Attachments = decodeStringSlice(attachmentsJSON)
	m.Reactions = decodeStringSlice(reactionsJSON)
	m.ReplyToID = replyTo.String
	return &m, nil
}

// AddReactionToMessage appends a reaction to a stored message if it is not already present.
func (r *Repository) AddReactionToMessage(channelID, messageID, reaction string) error {
	return r.mutateMessageReactions(channelID, messageID, func(reactions []string) []string {
		reaction = strings.TrimSpace(reaction)
		if reaction == "" {
			return reactions
		}
		for _, existing := range reactions {
			if existing == reaction {
				return reactions
			}
		}
		return append(reactions, reaction)
	})
}

// RemoveReactionFromMessage removes a reaction from a stored message.
func (r *Repository) RemoveReactionFromMessage(channelID, messageID, reaction string) error {
	return r.mutateMessageReactions(channelID, messageID, func(reactions []string) []string {
		reaction = strings.TrimSpace(reaction)
		if reaction == "" {
			return reactions
		}
		filtered := reactions[:0]
		for _, existing := range reactions {
			if existing != reaction {
				filtered = append(filtered, existing)
			}
		}
		return filtered
	})
}

// ClearReactionsForMessage clears all tracked reactions for a stored message.
func (r *Repository) ClearReactionsForMessage(channelID, messageID string) error {
	_, err := r.db.Exec(`UPDATE messages SET reactions_json = '[]' WHERE channel_id = ? AND discord_id = ?`, channelID, messageID)
	return err
}

func (r *Repository) mutateMessageReactions(channelID, messageID string, mutator func([]string) []string) error {
	tx, err := r.BeginImmediate()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var reactionsJSON string
	err = tx.QueryRow(`
		SELECT reactions_json
		FROM messages
		WHERE channel_id = ? AND discord_id = ?
	`, channelID, messageID).Scan(&reactionsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	reactions := decodeStringSlice(reactionsJSON)
	reactions = mutator(reactions)
	updatedJSON, err := json.Marshal(reactions)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE messages
		SET reactions_json = ?
		WHERE channel_id = ? AND discord_id = ?
	`, string(updatedJSON), channelID, messageID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetMessagesSinceLastBotReply fetches all messages since the bot's last response.
func (r *Repository) GetMessagesSinceLastBotReply(channelID string) ([]Message, error) {
	var lastBotTimestamp time.Time
	err := r.db.QueryRow(`
		SELECT timestamp FROM messages 
		WHERE channel_id = ? AND is_bot = 1 AND is_deleted = 0
		ORDER BY timestamp DESC LIMIT 1
	`, channelID).Scan(&lastBotTimestamp)

	if err != nil {
		if err == sql.ErrNoRows {
			// No bot reply yet, get last 50
			return r.GetRecentMessages(channelID, 50)
		}
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT id, discord_id, author, author_id, content, channel_id, category, mentions_json, attachments_json, reactions_json, is_bot, timestamp, reply_to_id, episode_id
		FROM messages
		WHERE channel_id = ? AND timestamp > ? AND is_deleted = 0
		ORDER BY timestamp ASC
	`, channelID, lastBotTimestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		var replyTo sql.NullString
		var mentionsJSON, attachmentsJSON, reactionsJSON string
		err := rows.Scan(&m.ID, &m.DiscordID, &m.Author, &m.AuthorID, &m.Content, &m.ChannelID, &m.Category, &mentionsJSON, &attachmentsJSON, &reactionsJSON, &m.IsBot, &m.Timestamp, &replyTo, &m.EpisodeID)
		if err != nil {
			return nil, err
		}
		m.Mentions = decodeStringSlice(mentionsJSON)
		m.Attachments = decodeStringSlice(attachmentsJSON)
		m.Reactions = decodeStringSlice(reactionsJSON)
		m.ReplyToID = replyTo.String
		msgs = append(msgs, m)
	}
	return msgs, nil
}
