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

type CategorySpec struct {
	Name     string
	Channels []string
}

func handleProvision5x5(s *discordgo.Session, guildID string) (string, error) {
	if s == nil {
		return "Sessão do Discord não iniciada.", nil
	}

	matrix := []CategorySpec{
		{
			Name: "🏛️ 1. QG & GOVERNANÇA",
			Channels: []string{
				"anuncios-oficiais",
				"qg-geral",
				"diretrizes-e-regras",
				"assembleia-e-pautas",
				"micelia-comandos",
			},
		},
		{
			Name: "📖 2. FORMAÇÃO & IDEOLOGIA",
			Channels: []string{
				"biblioteca-urupe",
				"debates-e-teoria",
				"agroecologia-e-terra",
				"historia-e-luta",
				"oficinas-e-cursos",
			},
		},
		{
			Name: "⚔️ 3. GUERRILHA DIGITAL & AGITPROP",
			Channels: []string{
				"central-de-brigadas",
				"design-e-memes",
				"audiovisual-e-cortes",
				"disparo-e-campanhas",
				"ciberseguranca",
			},
		},
		{
			Name: "🌐 4. REDES SOCIAIS & OBSERVAÇÃO",
			Channels: []string{
				"instagram-e-tiktok",
				"x-e-bluesky",
				"podcasts-e-midia",
				"observatorio-e-clipping",
				"metricas-e-impacto",
			},
		},
		{
			Name: "🌿 5. CULTURA & COMUNIDADE URUPÊ",
			Channels: []string{
				"ecos-do-micelio",
				"arte-e-literatura",
				"musica-e-playlists",
				"acolhimento-e-ajuda",
				"cafe-e-conversas",
			},
		},
	}

	createdCategories := 0
	createdChannels := 0

	for _, catSpec := range matrix {
		cat, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name: catSpec.Name,
			Type: discordgo.ChannelTypeGuildCategory,
		})
		if err != nil {
			return fmt.Sprintf("Erro ao criar categoria %s: %v", catSpec.Name, err), nil
		}
		createdCategories++

		for _, chName := range catSpec.Channels {
			_, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
				Name:     chName,
				Type:     discordgo.ChannelTypeGuildText,
				ParentID: cat.ID,
			})
			if err == nil {
				createdChannels++
			}
		}
	}

	return fmt.Sprintf(" Matriz Ontológica 5x5 da Rede Urupê provisionada com sucesso!\n- **%d Categorias** criadas\n- **%d Canais de Texto** estruturados.", createdCategories, createdChannels), nil
}
