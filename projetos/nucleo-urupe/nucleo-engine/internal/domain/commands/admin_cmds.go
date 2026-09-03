/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"nucleo-engine/internal/data/sqlite"
)

func AdminCommand(repo *sqlite.Repository) *Cmd {
	return &Cmd{
		Def: &discordgo.ApplicationCommand{
			Name:        "admin",
			Description: "Painel de administracao do Maze",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "welcome-set",
					Description: "Configurar mensagem de boas-vindas",
					Options: []*discordgo.ApplicationCommandOption{
						{Type: discordgo.ApplicationCommandOptionChannel, Name: "canal", Description: "Canal para boas-vindas", Required: true},
						{Type: discordgo.ApplicationCommandOptionString, Name: "mensagem", Description: "Mensagem (use {user} e {server})", Required: true},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "goodbye-set",
					Description: "Configurar mensagem de despedida",
					Options: []*discordgo.ApplicationCommandOption{
						{Type: discordgo.ApplicationCommandOptionChannel, Name: "canal", Description: "Canal para despedidas", Required: true},
						{Type: discordgo.ApplicationCommandOptionString, Name: "mensagem", Description: "Mensagem (use {user} e {server})", Required: true},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "welcome-toggle",
					Description: "Ativar/desativar boas-vindas",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "audit",
					Description: "Ver log de auditoria",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "stats",
					Description: "Estatisticas do servidor",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "provision-5x5",
					Description: "Provisionar a Matriz Ontologica 5x5 da Rede Urupe no Discord",
				},
			},
		},
		Handle: func(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
			data := i.ApplicationCommandData()
			sub := data.Options[0].Name

			switch sub {
			case "welcome-set":
				return handleWelcomeSet(repo, data, i.GuildID, s)
			case "goodbye-set":
				return handleGoodbyeSet(repo, data, i.GuildID, s)
			case "welcome-toggle":
				return handleWelcomeToggle(repo, i.GuildID)
			case "audit":
				return handleAuditLog(repo, i.GuildID)
			case "stats":
				return handleServerStats(repo, i.GuildID, i.ChannelID)
			case "provision-5x5":
				return handleProvision5x5(s, i.GuildID)
			default:
				return "Subcomando desconhecido.", nil
			}
		},
	}
}

func handleWelcomeSet(repo *sqlite.Repository, data discordgo.ApplicationCommandInteractionData, guildID string, s *discordgo.Session) (string, error) {
	var channelID, message string
	for _, opt := range data.Options[0].Options {
		switch opt.Name {
		case "canal":
			channelID = opt.ChannelValue(s).ID
		case "mensagem":
			message = opt.StringValue()
		}
	}

	_, err := repo.GetDB().Exec(`
		INSERT INTO welcome_config (guild_id, enabled, channel_id, welcome_message, updated_at)
		VALUES (?, 1, ?, ?, datetime('now'))
		ON CONFLICT(guild_id) DO UPDATE SET
			enabled = 1,
			channel_id = excluded.channel_id,
			welcome_message = excluded.welcome_message,
			updated_at = excluded.updated_at
	`, guildID, channelID, message)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Mensagem de boas-vindas configurada em <#%s>", channelID), nil
}

func handleGoodbyeSet(repo *sqlite.Repository, data discordgo.ApplicationCommandInteractionData, guildID string, s *discordgo.Session) (string, error) {
	var channelID, message string
	for _, opt := range data.Options[0].Options {
		switch opt.Name {
		case "canal":
			channelID = opt.ChannelValue(s).ID
		case "mensagem":
			message = opt.StringValue()
		}
	}

	_, err := repo.GetDB().Exec(`
		INSERT INTO welcome_config (guild_id, enabled, channel_id, goodbye_message, updated_at)
		VALUES (?, 1, ?, ?, datetime('now'))
		ON CONFLICT(guild_id) DO UPDATE SET
			enabled = 1,
			channel_id = excluded.channel_id,
			goodbye_message = excluded.goodbye_message,
			updated_at = excluded.updated_at
	`, guildID, channelID, message)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Mensagem de despedida configurada em <#%s>", channelID), nil
}

func handleWelcomeToggle(repo *sqlite.Repository, guildID string) (string, error) {
	var enabled bool
	err := repo.GetDB().QueryRow("SELECT enabled FROM welcome_config WHERE guild_id = ?", guildID).Scan(&enabled)
	if err != nil {
		return "Sistema de boas-vindas nao configurado. Use `/admin welcome-set` primeiro.", nil
	}

	newState := !enabled
	_, err = repo.GetDB().Exec("UPDATE welcome_config SET enabled = ?, updated_at = datetime('now') WHERE guild_id = ?", newState, guildID)
	if err != nil {
		return "", err
	}

	state := "desativado"
	if newState {
		state = "ativado"
	}
	return fmt.Sprintf("Sistema de boas-vindas: %s", state), nil
}

func handleAuditLog(repo *sqlite.Repository, guildID string) (string, error) {
	events, err := repo.GetAuditLog(guildID, 15)
	if err != nil || len(events) == 0 {
		return "Nenhum evento de auditoria.", nil
	}

	var b strings.Builder
	b.WriteString("**Log de Auditoria**\n")
	for _, e := range events {
		icon := "📋"
		switch e.Action {
		case "warn":
			icon = "⚠️"
		case "kick":
			icon = "👢"
		case "ban":
			icon = "🔨"
		case "delete":
			icon = "🗑️"
		case "timeout":
			icon = "⏰"
		case "join":
			icon = "➡️"
		case "leave":
			icon = "⬅️"
		}
		b.WriteString(fmt.Sprintf("%s **%s** - %s → %s (%s)\n", icon, e.Action, e.ActorName, e.TargetName, e.CreatedAt))
	}

	return b.String(), nil
}

func handleServerStats(repo *sqlite.Repository, guildID, channelID string) (string, error) {
	totalMembers, _ := repo.CountMemberProfiles()

	var msgCount int
	repo.GetDB().QueryRow("SELECT COUNT(*) FROM messages WHERE is_deleted = 0").Scan(&msgCount)

	var botMsgCount int
	repo.GetDB().QueryRow("SELECT COUNT(*) FROM messages WHERE is_bot = 1 AND is_deleted = 0").Scan(&botMsgCount)

	var warningCount int
	repo.GetDB().QueryRow("SELECT COUNT(*) FROM guild_warnings WHERE guild_id = ?", guildID).Scan(&warningCount)

	var eventCount int
	repo.GetDB().QueryRow("SELECT COUNT(*) FROM member_events WHERE guild_id = ?", guildID).Scan(&eventCount)

	totalTokens, _ := repo.GetTotalCostLast24Hours()

	var b strings.Builder
	b.WriteString("**Estatisticas do Servidor**\n")
	b.WriteString(fmt.Sprintf("- Membros registrados: %d\n", totalMembers))
	b.WriteString(fmt.Sprintf("- Total de mensagens: %d\n", msgCount))
	b.WriteString(fmt.Sprintf("- Mensagens do bot: %d\n", botMsgCount))
	b.WriteString(fmt.Sprintf("- Avisos ativos: %d\n", warningCount))
	b.WriteString(fmt.Sprintf("- Eventos de membros: %d\n", eventCount))
	b.WriteString(fmt.Sprintf("- Tokens (24h): %d\n", totalTokens))

	return b.String(), nil
}

type ChannelSpec struct {
	Name string
	Type discordgo.ChannelType
}

type CategorySpec struct {
	Name     string
	Channels []ChannelSpec
}

func handleProvision5x5(s *discordgo.Session, guildID string) (string, error) {
	if s == nil {
		return "Sessão do Discord não iniciada.", nil
	}

	topChannels := []ChannelSpec{
		{Name: "📜│leia-me", Type: discordgo.ChannelTypeGuildText},
	}

	matrix := []CategorySpec{
		{
			Name: "📌 INFORMAÇÕES",
			Channels: []ChannelSpec{
				{Name: "🌐│i-mundo", Type: discordgo.ChannelTypeGuildForum},
				{Name: "🏛️│ii-sociedade", Type: discordgo.ChannelTypeGuildForum},
				{Name: "⚙️│iii-cosmotecnica", Type: discordgo.ChannelTypeGuildForum},
				{Name: "⚔️│iv-praxis", Type: discordgo.ChannelTypeGuildForum},
				{Name: "🔥│v-espirito", Type: discordgo.ChannelTypeGuildForum},
			},
		},
		{
			Name: "💬 CANAIS COMUNS",
			Channels: []ChannelSpec{
				{Name: "📢│anuncios", Type: discordgo.ChannelTypeGuildText},
				{Name: "🌱│apresentacoes", Type: discordgo.ChannelTypeGuildText},
				{Name: "💬│chat-comum", Type: discordgo.ChannelTypeGuildText},
				{Name: "🎭│chat-memetico", Type: discordgo.ChannelTypeGuildText},
				{Name: "📐│chat-serio", Type: discordgo.ChannelTypeGuildText},
				{Name: "🍄│micorriza", Type: discordgo.ChannelTypeGuildText},
			},
		},
		{
			Name: "🎨 CANAIS CULTURAIS",
			Channels: []ChannelSpec{
				{Name: "📰│noticias", Type: discordgo.ChannelTypeGuildText},
				{Name: "📹│videos", Type: discordgo.ChannelTypeGuildText},
				{Name: "🎮│jogos", Type: discordgo.ChannelTypeGuildText},
				{Name: "🎬│cinema", Type: discordgo.ChannelTypeGuildText},
				{Name: "🎵│musica", Type: discordgo.ChannelTypeGuildText},
				{Name: "📚│literatura", Type: discordgo.ChannelTypeGuildText},
			},
		},
		{
			Name: "📖 FRENTES DE ESTUDO",
			Channels: []ChannelSpec{
				{Name: "🏛️│filosofia", Type: discordgo.ChannelTypeGuildText},
				{Name: "🚩│politica", Type: discordgo.ChannelTypeGuildText},
				{Name: "🎓│pedagogia", Type: discordgo.ChannelTypeGuildText},
				{Name: "🌿│ecologia", Type: discordgo.ChannelTypeGuildText},
				{Name: "⏳│historia", Type: discordgo.ChannelTypeGuildText},
				{Name: "🧠│humanidades", Type: discordgo.ChannelTypeGuildText},
				{Name: "🪞│psicologia", Type: discordgo.ChannelTypeGuildText},
				{Name: "💻│tecnologia", Type: discordgo.ChannelTypeGuildText},
				{Name: "🔧│engenharia", Type: discordgo.ChannelTypeGuildText},
			},
		},
		{
			Name: "⚔️ FRENTES DE ATUAÇÃO",
			Channels: []ChannelSpec{
				{Name: "🍄│nucleo-urupe", Type: discordgo.ChannelTypeGuildText},
				{Name: "🌾│rizoma", Type: discordgo.ChannelTypeGuildText},
				{Name: "📱│app-urupe", Type: discordgo.ChannelTypeGuildText},
				{Name: "🧫│spore-ops", Type: discordgo.ChannelTypeGuildText},
				{Name: "🐝│jatai-ops", Type: discordgo.ChannelTypeGuildText},
			},
		},
	}

	createdCategories := 0
	createdChannels := 0

	// 1. Criar canal raiz (Topo)
	for _, chSpec := range topChannels {
		_, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name: chSpec.Name,
			Type: chSpec.Type,
		})
		if err == nil {
			createdChannels++
		}
	}

	// 2. Criar categorias e seus respectivos canais/fóruns
	for _, catSpec := range matrix {
		cat, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name: catSpec.Name,
			Type: discordgo.ChannelTypeGuildCategory,
		})
		if err != nil {
			return fmt.Sprintf("Erro ao criar categoria %s: %v", catSpec.Name, err), nil
		}
		createdCategories++

		for _, chSpec := range catSpec.Channels {
			_, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
				Name:     chSpec.Name,
				Type:     chSpec.Type,
				ParentID: cat.ID,
			})
			if err == nil {
				createdChannels++
			}
		}
	}

	return fmt.Sprintf(" Matriz do Discord da Frente Urupê provisionada com sucesso!\n- **1 Canal Raiz** (#📜│leia-me)\n- **%d Categorias** criadas\n- **%d Canais & Fóruns** estruturados (no formato #emoji│nome-do-canal, com os 5 Fóruns do Manifesto e zero voz).", createdCategories, createdChannels), nil
}
