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
		{Name: "leia-me-manifesto-e-regras", Type: discordgo.ChannelTypeGuildText},
		{Name: "quarentena-solicitar-acesso", Type: discordgo.ChannelTypeGuildText},
	}

	matrix := []CategorySpec{
		{
			Name: "📌 1. INFORMAÇÕES & FÓRUNS",
			Channels: []ChannelSpec{
				{Name: "formacao-e-teoria-ecossocialista", Type: discordgo.ChannelTypeGuildForum},
				{Name: "acao-direta-e-brigadas", Type: discordgo.ChannelTypeGuildForum},
				{Name: "soberania-tecnologica-e-rizoma", Type: discordgo.ChannelTypeGuildForum},
				{Name: "arte-cultura-e-memetica", Type: discordgo.ChannelTypeGuildForum},
				{Name: "agroecologia-e-territorio", Type: discordgo.ChannelTypeGuildForum},
			},
		},
		{
			Name: "💬 2. CANAIS COMUNS",
			Channels: []ChannelSpec{
				{Name: "bate-papo-geral", Type: discordgo.ChannelTypeGuildText},
				{Name: "midia-e-compartilhamento", Type: discordgo.ChannelTypeGuildText},
				{Name: "🔊 Praça Pública (Voz)", Type: discordgo.ChannelTypeGuildVoice},
			},
		},
		{
			Name: "📖 3. FRENTES DE ESTUDO",
			Channels: []ChannelSpec{
				{Name: "clube-de-leitura", Type: discordgo.ChannelTypeGuildText},
				{Name: "debates-ideologicos", Type: discordgo.ChannelTypeGuildText},
				{Name: "🔊 Sala de Estudos (Voz)", Type: discordgo.ChannelTypeGuildVoice},
			},
		},
		{
			Name: "⚔️ 4. FRENTES DE ATUAÇÃO",
			Channels: []ChannelSpec{
				{Name: "brigadas-e-acao-direta", Type: discordgo.ChannelTypeGuildText},
				{Name: "agitprop-e-comunicacao", Type: discordgo.ChannelTypeGuildText},
				{Name: "imprensa-spore-ops", Type: discordgo.ChannelTypeGuildText},
			},
		},
		{
			Name: "🏛️ 5. INSTITUCIONAL E COMUNICAÇÃO",
			Channels: []ChannelSpec{
				{Name: "anuncios-oficiais", Type: discordgo.ChannelTypeGuildText},
				{Name: "fale-com-a-micelia", Type: discordgo.ChannelTypeGuildText},
				{Name: "memorias-da-frente", Type: discordgo.ChannelTypeGuildText},
			},
		},
	}

	createdCategories := 0
	createdChannels := 0

	// 1. Criar canais fora de categoria (Topo)
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

	return fmt.Sprintf(" Matriz Enxuta da Rede Urupê provisionada com sucesso no Discord!\n- **2 Canais de Topo** (#leia-me e #quarentena)\n- **%d Categorias** criadas\n- **%d Canais & Fóruns** estruturados (com 5 Fóruns Temáticos e Canais de Estudo/Atuação/Comunicação).", createdCategories, createdChannels), nil
}
