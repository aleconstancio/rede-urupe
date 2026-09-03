/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"nucleo-engine/internal/data/sqlite"
	"github.com/aleconstancio/talos/v2/engine/intelligence"
	minotaur "github.com/aleconstancio/talos/v2/behavior/persona"
)

type DiscordMessenger struct {
	session     *discordgo.Session
	repo        *sqlite.Repository
	broadcaster Broadcaster
	emojis      map[string]string // Name -> ID
}

func NewDiscordMessenger(session *discordgo.Session, repo *sqlite.Repository, broadcaster Broadcaster) *DiscordMessenger {
	return &DiscordMessenger{
		session:     session,
		repo:        repo,
		broadcaster: broadcaster,
		emojis:      make(map[string]string),
	}
}

func (m *DiscordMessenger) SetEmojis(emojis map[string]string) {
	m.emojis = emojis
}

func (m *DiscordMessenger) SendAndStoreReply(channelID, replyToID, content, monologue, ledger string, reactions []string) error {
	if m.session == nil {
		return fmt.Errorf("discord session not configured")
	}

	// Resolve reactions
	if replyToID != "" && len(reactions) > 0 {
		go func() {
			for _, emoji := range reactions {
				resolved := emoji
				if id, ok := m.emojis[emoji]; ok {
					resolved = fmt.Sprintf("%s:%s", emoji, id)
				}
				_ = m.session.MessageReactionAdd(channelID, replyToID, resolved)
				time.Sleep(200 * time.Millisecond)
			}
		}()
	}

	if strings.TrimSpace(content) == "" {
		return nil
	}

	// Auto-format custom emojis in content: :name: -> <:name:id>
	for name, id := range m.emojis {
		content = strings.ReplaceAll(content, ":"+name+":", fmt.Sprintf("<:%s:%s>", name, id))
	}

	if len(content) > 2000 {
		content = content[:1997] + "..."
	}


	msg, err := m.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		Reference: &discordgo.MessageReference{
			MessageID: replyToID,
			ChannelID: channelID,
		},
	})
	if err != nil {
		return err
	}

	authorID := ""
	if msg.Author != nil {
		authorID = msg.Author.ID
	}

	_, err = m.repo.SaveMessage(
		msg.ID,
		MessageAuthorName(msg),
		authorID,
		msg.Content,
		msg.ChannelID,
		ClassifyMessageCategory(msg.Content),
		true,
		time.Time(msg.Timestamp),
		replyToID,
		nil, nil, nil,
		monologue,
		ledger,
	)
	return err
}

func (m *DiscordMessenger) SyncDiscordProfile(ctx context.Context, channelID string, identity minotaur.CoreIdentityProfile) {
	if m.session == nil || identity.DisplayName == "" {
		return
	}

	ch, err := m.session.State.Channel(channelID)
	if err != nil {
		ch, err = m.session.Channel(channelID)
	}
	if err == nil && ch.GuildID != "" {
		member, err := m.session.State.Member(ch.GuildID, m.session.State.User.ID)
		if err != nil {
			member, _ = m.session.GuildMember(ch.GuildID, m.session.State.User.ID)
		}

		if member != nil && member.Nick != identity.DisplayName {
			log.Printf("[Messenger] Updating Discord nickname to %s", identity.DisplayName)
			_ = m.session.GuildMemberNickname(ch.GuildID, "@me", identity.DisplayName)
		}
	}
}

func (m *DiscordMessenger) Broadcast(event string) {
	if m.broadcaster != nil {
		m.broadcaster.Broadcast(event)
	}
}


func MessageAuthorName(m *discordgo.Message) string {
	if m == nil || m.Author == nil {
		return "Bot"
	}
	if m.Member != nil && m.Member.Nick != "" {
		return m.Member.Nick
	}
	return m.Author.Username
}

func ClassifyMessageCategory(content string) string {
	return string(intelligence.ClassifyDetailed(content).Primary)
}
