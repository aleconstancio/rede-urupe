/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package discord

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/aleconstancio/talos/v2/engine/intelligence"
	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/domain/commands"
	"nucleo-engine/internal/domain/gateway"
	"nucleo-engine/internal/domain/moderation"
	"nucleo-engine/internal/pkg/timeutil"
)

type Handler struct {
	repo            *sqlite.Repository
	targetGuildID   string
	targetChannelID string
	triggerChan     chan<- gateway.TriggerEvent
	cmdRouter       *commands.Router
	enforcer        *moderation.Enforcer
	session         *discordgo.Session
	gatewayWorker   *gateway.GatewayWorker
	multiChannel    bool

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewHandler(triggerChan chan<- gateway.TriggerEvent, repo *sqlite.Repository, targetGuildID, targetChannelID string) *Handler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Handler{
		repo:            repo,
		targetGuildID:   targetGuildID,
		targetChannelID: targetChannelID,
		triggerChan:     triggerChan,
		enforcer:        moderation.NewEnforcer(repo),
		ctx:             ctx,
		cancel:          cancel,
	}
}

func (h *Handler) SetGatewayWorker(w *gateway.GatewayWorker) {
	h.gatewayWorker = w
}

func (h *Handler) SetMultiChannel(multi bool) {
	h.multiChannel = multi
}

func (h *Handler) SetCommandRouter(r *commands.Router) {
	h.cmdRouter = r
}

func (h *Handler) Start(s *discordgo.Session) {
	h.session = s
	h.enforcer.SetSession(moderation.NewDiscordSession(s))
	log.Printf("[Discord] Bot started with ID: %s", s.State.User.ID)

	if h.cmdRouter != nil {
		cmds := h.cmdRouter.AllCommands()
		for _, cmd := range cmds {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
			if err != nil {
				log.Printf("[Discord] Failed to register global command %q: %v", cmd.Name, err)
			} else {
				log.Printf("[Discord] Registered global command: %s", cmd.Name)
			}
		}
	}
}

func (h *Handler) Close() {
	log.Println("[Discord] Shutting down handler...")
	h.cancel()
	h.wg.Wait()
}

func (h *Handler) OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil || m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}
	if h.targetGuildID != "" && m.GuildID != h.targetGuildID {
		return
	}
	if !h.multiChannel && m.ChannelID != h.targetChannelID {
		return
	}

	log.Printf("[Discord] Received message from %s: %s", m.Author.Username, m.Content)

	if err := h.ensureGuildConfig(m.GuildID); err != nil {
		log.Printf("[Discord] Failed to ensure guild config: %v", err)
	}

	rowID, err := h.saveMessage(m.Message, false)
	if err != nil {
		log.Printf("Error saving message: %v", err)
		return
	}

	h.runModeration(m)

	if h.multiChannel && h.gatewayWorker != nil {
		h.gatewayWorker.AddChannel(m.ChannelID)
		_ = h.repo.RegisterChannel(m.GuildID, m.ChannelID)
	}

	if h.triggerChan == nil {
		log.Printf("[Discord] Trigger channel not configured")
		return
	}

	select {
	case h.triggerChan <- gateway.TriggerEvent{
		ChannelID:    m.ChannelID,
		MessageRowID: rowID,
		MessageID:    m.ID,
		Timestamp:    timeutil.InBrasilia(time.Time(m.Timestamp)),
		IsReactive:   h.isDirectlyAddressed(m, s.State.User.ID),
	}:
	default:
		log.Printf("[CRITICAL] GatewayWorker trigger channel full, dropping event for message %s", m.ID)
	}
}

func (h *Handler) runModeration(m *discordgo.MessageCreate) {
	if h.enforcer == nil || h.targetGuildID == "" {
		return
	}

	gc, err := h.repo.GetGuildConfig(m.GuildID)
	if err != nil || !gc.ModEnabled {
		return
	}

	tone := intelligence.AnalyzeTone(strings.ToLower(m.Content), m.Content)

	var recentContent []string
	recent, _ := h.repo.GetRecentMessages(m.ChannelID, 10)
	for _, msg := range recent {
		recentContent = append(recentContent, msg.Content)
	}

	ctx := moderation.ModContext{
		Content:   m.Content,
		AuthorID:  m.Author.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
		MessageID: m.ID,
		Recent:    recentContent,
	}

	hits := moderation.CheckMessage(ctx, tone)
	if len(hits) > 0 {
		h.enforcer.ProcessHits(hits, ctx)
	}
}

func (h *Handler) saveMessage(m *discordgo.Message, isBot bool) (int64, error) {
	var mentions []string
	for _, user := range m.Mentions {
		mentions = append(mentions, user.ID)
	}
	var attachments []string
	for _, a := range m.Attachments {
		attachments = append(attachments, a.URL)
	}
	var reactions []string
	for _, r := range m.Reactions {
		reactions = append(reactions, r.Emoji.Name)
	}

	var roles []string
	if m.Member != nil {
		for _, roleID := range m.Member.Roles {
			// Extract role IDs. We will resolve names in the worker if needed.
			roles = append(roles, roleID)
		}
	}

	// Update persistent member profile
	_ = h.repo.UpsertMemberProfile(&sqlite.MemberProfile{
		DiscordID: authorID(m),
		Name:      messageAuthorName(m),
		Roles:     roles,
	})

	return h.repo.SaveMessage(m.ID, messageAuthorName(m), authorID(m), m.Content, m.ChannelID, classifyMessageCategory(m.Content), isBot, time.Time(m.Timestamp), replyToID(m), mentions, attachments, reactions, "", "")
}

func (h *Handler) isDirectlyAddressed(m *discordgo.MessageCreate, botID string) bool {
	if m == nil || botID == "" {
		return false
	}

	// 1. Reagir sempre se a mensagem for enviada no canal #micorriza ou #fale-com-a-micelia
	if h.session != nil {
		ch, err := h.session.State.Channel(m.ChannelID)
		if err == nil && ch != nil {
			chName := strings.ToLower(ch.Name)
			if strings.Contains(chName, "micorriza") || strings.Contains(chName, "micelia") {
				return true
			}
		}
	}

	if strings.Contains(m.Content, "<@"+botID+">") || strings.Contains(m.Content, "<@!"+botID+">") {
		return true
	}
	for _, user := range m.Mentions {
		if user.ID == botID {
			return true
		}
	}
	if m.ReferencedMessage != nil && m.ReferencedMessage.Author != nil && m.ReferencedMessage.Author.ID == botID {
		return true
	}
	c := strings.ToLower(m.Content)
	identities, _ := h.repo.ListCoreIdentityProfiles()
	for _, id := range identities {
		if id.IsEnabled && strings.Contains(c, strings.ToLower(id.DisplayName)) {
			return true
		}
	}
	return false
}

func classifyMessageCategory(content string) string {
	return string(intelligence.ClassifyDetailed(content).Primary)
}

func messageAuthorName(m *discordgo.Message) string {
	if m == nil || m.Author == nil {
		return "Unknown"
	}
	if m.Member != nil && m.Member.Nick != "" {
		return m.Member.Nick
	}
	return m.Author.Username
}

func authorID(m *discordgo.Message) string {
	if m == nil || m.Author == nil {
		return ""
	}
	return m.Author.ID
}

func replyToID(m *discordgo.Message) string {
	if m != nil && m.ReferencedMessage != nil {
		return m.ReferencedMessage.ID
	}
	return ""
}

func (h *Handler) OnMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author == nil || m.Author.ID == s.State.User.ID || m.Author.Bot {
		return
	}
	if !h.multiChannel && m.ChannelID != h.targetChannelID {
		return
	}

	// For updates, we often get a partial message.
	// If content is empty, it might be an embed update or something we don't care about for context.
	if m.Content == "" && len(m.Attachments) == 0 {
		return
	}

	_, err := h.saveMessage(m.Message, false)
	if err != nil {
		log.Printf("Error updating message: %v", err)
	}
}

func (h *Handler) OnMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	if m == nil || (!h.multiChannel && m.ChannelID != h.targetChannelID) {
		return
	}
	if err := h.repo.MarkMessageDeleted(m.ChannelID, m.ID, timeutil.Now()); err != nil {
		log.Printf("Error marking deleted message: %v", err)
	}
}

func (h *Handler) OnReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m == nil || (!h.multiChannel && m.ChannelID != h.targetChannelID) {
		return
	}
	if err := h.repo.AddReactionToMessage(m.ChannelID, m.MessageID, m.Emoji.Name); err != nil {
		log.Printf("Error updating reaction add: %v", err)
	}
}

func (h *Handler) OnReactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	if m == nil || (!h.multiChannel && m.ChannelID != h.targetChannelID) {
		return
	}
	if err := h.repo.RemoveReactionFromMessage(m.ChannelID, m.MessageID, m.Emoji.Name); err != nil {
		log.Printf("Error updating reaction remove: %v", err)
	}
}

func (h *Handler) OnReactionRemoveAll(s *discordgo.Session, m *discordgo.MessageReactionRemoveAll) {
	if m == nil || (!h.multiChannel && m.ChannelID != h.targetChannelID) {
		return
	}
	if err := h.repo.ClearReactionsForMessage(m.ChannelID, m.MessageID); err != nil {
		log.Printf("Error clearing reactions: %v", err)
	}
}

func (h *Handler) OnGuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if m.GuildID != h.targetGuildID {
		return
	}

	userName := m.User.Username
	if m.Nick != "" {
		userName = m.Nick
	}

	if err := h.repo.LogMemberEvent(m.GuildID, m.User.ID, userName, "join", "", ""); err != nil {
		log.Printf("[Discord] Failed to log member join: %v", err)
	}

	_ = h.repo.UpsertMemberProfile(&sqlite.MemberProfile{
		DiscordID: m.User.ID,
		Name:      userName,
		Roles:     m.Roles,
	})

	_ = h.repo.LogAuditEvent(m.GuildID, "member", "join", m.User.ID, userName, m.User.ID, userName, "", "")

	gc, err := h.repo.GetGuildConfig(m.GuildID)
	if err == nil && gc.ModLogChannelID != "" {
		msg := fmt.Sprintf("**[MEMBRO] Entrou no servidor**\n- Usuario: <@%s>\n- ID: `%s`", m.User.ID, m.User.ID)
		_, _ = h.repo.GetDB().Exec("INSERT INTO messages (discord_id, author, author_id, content, channel_id, category, is_bot, timestamp) VALUES (?, 'Maze', '0', ?, ?, 'system', 1, datetime('now'))",
			fmt.Sprintf("member-join-%s", m.User.ID), msg, gc.ModLogChannelID)
	}
	// Aprovação no canal #apresentacoes é realizada 100% manualmente por administradores.
}

func (h *Handler) OnGuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	if m.GuildID != h.targetGuildID {
		return
	}

	userName := m.User.Username
	if m.Nick != "" {
		userName = m.Nick
	}

	if err := h.repo.LogMemberEvent(m.GuildID, m.User.ID, userName, "leave", "", ""); err != nil {
		log.Printf("[Discord] Failed to log member leave: %v", err)
	}

	_ = h.repo.LogAuditEvent(m.GuildID, "member", "leave", m.User.ID, userName, m.User.ID, userName, "", "")

	gc, err := h.repo.GetGuildConfig(m.GuildID)
	if err == nil && gc.ModLogChannelID != "" {
		msg := fmt.Sprintf("**[MEMBRO] Saiu do servidor**\n- Usuario: %s\n- ID: `%s`", userName, m.User.ID)
		_, _ = h.repo.GetDB().Exec("INSERT INTO messages (discord_id, author, author_id, content, channel_id, category, is_bot, timestamp) VALUES (?, 'Maze', '0', ?, ?, 'system', 1, datetime('now'))",
			fmt.Sprintf("member-leave-%s", m.User.ID), msg, gc.ModLogChannelID)
	}

	h.sendGoodbyeMessage(m)
}

func (h *Handler) sendWelcomeMessage(m *discordgo.GuildMemberAdd) {
	if h.session == nil {
		return
	}

	var wc struct {
		Enabled        bool
		ChannelID      string
		WelcomeMessage string
	}
	err := h.repo.GetDB().QueryRow(`
		SELECT enabled, channel_id, welcome_message FROM welcome_config WHERE guild_id = ?
	`, m.GuildID).Scan(&wc.Enabled, &wc.ChannelID, &wc.WelcomeMessage)
	if err != nil || !wc.Enabled || wc.ChannelID == "" {
		return
	}

	msg := wc.WelcomeMessage
	if msg == "" {
		msg = fmt.Sprintf("Bem-vindo(a) ao servidor, <@%s>!", m.User.ID)
	} else {
		msg = strings.ReplaceAll(msg, "{user}", fmt.Sprintf("<@%s>", m.User.ID))
		msg = strings.ReplaceAll(msg, "{server}", m.GuildID)
	}

	_, _ = h.session.ChannelMessageSend(wc.ChannelID, msg)
}

func (h *Handler) sendGoodbyeMessage(m *discordgo.GuildMemberRemove) {
	if h.session == nil {
		return
	}

	var wc struct {
		Enabled        bool
		ChannelID      string
		GoodbyeMessage string
	}
	err := h.repo.GetDB().QueryRow(`
		SELECT enabled, channel_id, goodbye_message FROM welcome_config WHERE guild_id = ?
	`, m.GuildID).Scan(&wc.Enabled, &wc.ChannelID, &wc.GoodbyeMessage)
	if err != nil || !wc.Enabled || wc.ChannelID == "" {
		return
	}

	msg := wc.GoodbyeMessage
	if msg == "" {
		msg = fmt.Sprintf("%s saiu do servidor.", m.User.Username)
	} else {
		msg = strings.ReplaceAll(msg, "{user}", m.User.Username)
		msg = strings.ReplaceAll(msg, "{server}", m.GuildID)
	}

	_, _ = h.session.ChannelMessageSend(wc.ChannelID, msg)
}

func (h *Handler) OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	if err := h.ensureGuildConfig(i.GuildID); err != nil {
		log.Printf("[Discord] Failed to ensure guild config: %v", err)
	}

	if h.cmdRouter != nil {
		h.cmdRouter.Handle(s, i)
	}
}

func (h *Handler) ensureGuildConfig(guildID string) error {
	_, err := h.repo.GetGuildConfig(guildID)
	if err != nil {
		return h.repo.UpsertGuildConfig(sqlite.GuildConfig{
			GuildID:     guildID,
			Prefix:      "/",
			ModEnabled:  false,
			MaxWarnings: 3,
		})
	}
	return nil
}
