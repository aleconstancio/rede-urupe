/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package moderation

import (
	"fmt"
	"log"
	"time"

	"nucleo-engine/internal/data/sqlite"
)

type Enforcer struct {
	repo    *sqlite.Repository
	session interface {
		GuildMemberRoleAdd(guildID, userID, roleID string) error
		ChannelMessageSend(channelID string, content string) (interface{}, error)
		ChannelMessageDelete(channelID, messageID string) error
		GuildMemberTimeout(guildID, userID string, until *time.Time) error
		GuildMemberRemove(guildID, userID, reason string) error
		GuildMemberBan(guildID, userID, reason string) error
	}
}

func NewEnforcer(repo *sqlite.Repository) *Enforcer {
	return &Enforcer{repo: repo}
}

func (e *Enforcer) SetSession(s interface{}) {
	if sess, ok := s.(interface {
		GuildMemberRoleAdd(guildID, userID, roleID string) error
		ChannelMessageSend(channelID string, content string) (interface{}, error)
		ChannelMessageDelete(channelID, messageID string) error
		GuildMemberTimeout(guildID, userID string, until *time.Time) error
		GuildMemberRemove(guildID, userID, reason string) error
		GuildMemberBan(guildID, userID, reason string) error
	}); ok {
		e.session = sess
	}
}

func (e *Enforcer) ProcessHits(hits []RuleHit, ctx ModContext) {
	if len(hits) == 0 {
		return
	}

	gc, err := e.repo.GetGuildConfig(ctx.GuildID)
	if err != nil || !gc.ModEnabled {
		return
	}

	warningCount, _ := e.repo.GetWarningCount(ctx.GuildID, ctx.AuthorID)

	for i := range hits {
		hit := &hits[i]

		if warningCount >= 5 && hit.Action == "warn" {
			hit.Action = "timeout"
			hit.ActionDuration = 15 * time.Minute
		} else if warningCount >= 3 && hit.Action == "warn" {
			hit.Action = "delete"
		}

		switch hit.Action {
		case "warn":
			e.enforceWarn(*hit, ctx, gc)
		case "delete":
			e.enforceDelete(*hit, ctx, gc)
		case "timeout":
			e.enforceTimeout(*hit, ctx, gc)
		}
	}
}

func (e *Enforcer) enforceWarn(hit RuleHit, ctx ModContext, gc *sqlite.GuildConfig) {
	severity := 1
	if hit.Score >= 0.9 {
		severity = 2
	}

	err := e.repo.AddWarning(ctx.GuildID, ctx.AuthorID, ctx.ChannelID, hit.Reason, "", severity)
	if err != nil {
		log.Printf("[Moderation] Failed to add warning: %v", err)
		return
	}

	count, _ := e.repo.GetWarningCount(ctx.GuildID, ctx.AuthorID)
	log.Printf("[Moderation] Warning issued to %s: %s (total: %d/%d)", ctx.AuthorID, hit.Reason, count, gc.MaxWarnings)

	_ = e.repo.LogAuditEvent(ctx.GuildID, "moderation", "warn", "", "Maze", ctx.AuthorID, ctx.AuthorID,
		sqlite.AuditDetailsJSON(map[string]interface{}{
			"rule":          hit.Rule,
			"reason":        hit.Reason,
			"score":         hit.Score,
			"warning_count": count,
			"max_warnings":  gc.MaxWarnings,
		}), ctx.ChannelID)

	if gc.ModLogChannelID != "" {
		msg := fmt.Sprintf("**[MOD] Aviso**\n- Usuario: <@%s>\n- Regra: `%s`\n- Motivo: %s\n- Avisos: %d/%d",
			ctx.AuthorID, hit.Rule, hit.Reason, count, gc.MaxWarnings)
		_, _ = e.repo.GetDB().Exec("INSERT INTO messages (discord_id, author, author_id, content, channel_id, category, is_bot, timestamp) VALUES (?, 'Maze', '0', ?, ?, 'moderation', 1, datetime('now'))",
			fmt.Sprintf("mod-%d", count), msg, gc.ModLogChannelID)
	}
}

func (e *Enforcer) enforceDelete(hit RuleHit, ctx ModContext, gc *sqlite.GuildConfig) {
	log.Printf("[Moderation] Delete requested for message by %s: %s", ctx.AuthorID, hit.Reason)

	if e.session != nil {
		if err := e.session.ChannelMessageDelete(ctx.ChannelID, ctx.MessageID); err != nil {
			log.Printf("[Moderation] Failed to delete message: %v", err)
		}
	}

	_ = e.repo.LogAuditEvent(ctx.GuildID, "moderation", "delete", "", "Maze", ctx.AuthorID, ctx.AuthorID,
		sqlite.AuditDetailsJSON(map[string]interface{}{
			"rule":   hit.Rule,
			"reason": hit.Reason,
			"score":  hit.Score,
		}), ctx.ChannelID)

	if gc.ModLogChannelID != "" {
		msg := fmt.Sprintf("**[MOD] Mensagem removida**\n- Usuario: <@%s>\n- Regra: `%s`\n- Motivo: %s",
			ctx.AuthorID, hit.Rule, hit.Reason)
		_, _ = e.repo.GetDB().Exec("INSERT INTO messages (discord_id, author, author_id, content, channel_id, category, is_bot, timestamp) VALUES (?, 'Maze', '0', ?, ?, 'moderation', 1, datetime('now'))",
			fmt.Sprintf("mod-del-%s", ctx.AuthorID), msg, gc.ModLogChannelID)
	}
}

func (e *Enforcer) enforceTimeout(hit RuleHit, ctx ModContext, gc *sqlite.GuildConfig) {
	log.Printf("[Moderation] Timeout requested for %s: %s (duration: %v)", ctx.AuthorID, hit.Reason, hit.ActionDuration)

	if e.session != nil && hit.ActionDuration > 0 {
		until := time.Now().Add(hit.ActionDuration)
		if err := e.session.GuildMemberTimeout(ctx.GuildID, ctx.AuthorID, &until); err != nil {
			log.Printf("[Moderation] Failed to timeout user: %v", err)
		}
	}

	_ = e.repo.LogAuditEvent(ctx.GuildID, "moderation", "timeout", "", "Maze", ctx.AuthorID, ctx.AuthorID,
		sqlite.AuditDetailsJSON(map[string]interface{}{
			"rule":     hit.Rule,
			"reason":   hit.Reason,
			"score":    hit.Score,
			"duration": hit.ActionDuration.String(),
		}), ctx.ChannelID)

	if gc.ModLogChannelID != "" {
		msg := fmt.Sprintf("**[MOD] Timeout**\n- Usuario: <@%s>\n- Regra: `%s`\n- Motivo: %s\n- Duracao: %v",
			ctx.AuthorID, hit.Rule, hit.Reason, hit.ActionDuration)
		_, _ = e.repo.GetDB().Exec("INSERT INTO messages (discord_id, author, author_id, content, channel_id, category, is_bot, timestamp) VALUES (?, 'Maze', '0', ?, ?, 'moderation', 1, datetime('now'))",
			fmt.Sprintf("mod-timeout-%s", ctx.AuthorID), msg, gc.ModLogChannelID)
	}
}

func (e *Enforcer) EnforceManualAction(action, actorID, actorName, targetID, targetName, reason, guildID, channelID string) error {
	gc, err := e.repo.GetGuildConfig(guildID)
	if err != nil {
		return err
	}

	switch action {
	case "kick":
		if e.session != nil {
			if err := e.session.GuildMemberRemove(guildID, targetID, reason); err != nil {
				log.Printf("[Moderation] Failed to kick user %s: %v", targetID, err)
			}
		}
		_ = e.repo.LogAuditEvent(guildID, "moderation", "kick", actorID, actorName, targetID, targetName,
			sqlite.AuditDetailsJSON(map[string]interface{}{"reason": reason}), channelID)
		_ = e.repo.LogMemberEvent(guildID, targetID, targetName, "kick", reason, actorID)

	case "ban":
		if e.session != nil {
			if err := e.session.GuildMemberBan(guildID, targetID, reason); err != nil {
				log.Printf("[Moderation] Failed to ban user %s: %v", targetID, err)
			}
		}
		_ = e.repo.LogAuditEvent(guildID, "moderation", "ban", actorID, actorName, targetID, targetName,
			sqlite.AuditDetailsJSON(map[string]interface{}{"reason": reason}), channelID)
		_ = e.repo.LogMemberEvent(guildID, targetID, targetName, "ban", reason, actorID)

	case "warn":
		err = e.repo.AddWarning(guildID, targetID, channelID, reason, "", 1)
		if err != nil {
			return err
		}
		count, _ := e.repo.GetWarningCount(guildID, targetID)
		_ = e.repo.LogAuditEvent(guildID, "moderation", "warn", actorID, actorName, targetID, targetName,
			sqlite.AuditDetailsJSON(map[string]interface{}{"reason": reason, "warning_count": count, "max_warnings": gc.MaxWarnings}), channelID)

	case "clear_warnings":
		_, err = e.repo.GetDB().Exec("DELETE FROM guild_warnings WHERE guild_id = ? AND user_id = ?", guildID, targetID)
		_ = e.repo.LogAuditEvent(guildID, "moderation", "clear_warnings", actorID, actorName, targetID, targetName,
			sqlite.AuditDetailsJSON(map[string]interface{}{"reason": reason}), channelID)
	}

	return nil
}
