/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package moderation

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type DiscordSession struct {
	session *discordgo.Session
}

func NewDiscordSession(s *discordgo.Session) *DiscordSession {
	return &DiscordSession{session: s}
}

func (d *DiscordSession) GuildMemberRoleAdd(guildID, userID, roleID string) error {
	return d.session.GuildMemberRoleAdd(guildID, userID, roleID)
}

func (d *DiscordSession) ChannelMessageSend(channelID, content string) (interface{}, error) {
	return d.session.ChannelMessageSend(channelID, content)
}

func (d *DiscordSession) ChannelMessageDelete(channelID, messageID string) error {
	return d.session.ChannelMessageDelete(channelID, messageID)
}

func (d *DiscordSession) GuildMemberTimeout(guildID, userID string, until *time.Time) error {
	return d.session.GuildMemberTimeout(guildID, userID, until)
}

func (d *DiscordSession) GuildMemberRemove(guildID, userID, reason string) error {
	return d.session.GuildMemberDeleteWithReason(guildID, userID, reason)
}

func (d *DiscordSession) GuildMemberBan(guildID, userID, reason string) error {
	return d.session.GuildBanCreateWithReason(guildID, userID, reason, 0)
}
