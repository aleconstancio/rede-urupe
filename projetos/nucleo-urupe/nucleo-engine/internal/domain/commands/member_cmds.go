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

func MemberCommand(repo *sqlite.Repository) *Cmd {
	return &Cmd{
		Def: &discordgo.ApplicationCommand{
			Name:        "member",
			Description: "Gerenciar membros do servidor",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "info",
					Description: "Ver informacoes de um membro",
					Options: []*discordgo.ApplicationCommandOption{
						{Type: discordgo.ApplicationCommandOptionUser, Name: "usuario", Description: "Membro para consultar", Required: true},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "warn",
					Description: "Avisar um membro",
					Options: []*discordgo.ApplicationCommandOption{
						{Type: discordgo.ApplicationCommandOptionUser, Name: "usuario", Description: "Membro para avisar", Required: true},
						{Type: discordgo.ApplicationCommandOptionString, Name: "motivo", Description: "Motivo do aviso", Required: true},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "warnings",
					Description: "Ver avisos de um membro",
					Options: []*discordgo.ApplicationCommandOption{
						{Type: discordgo.ApplicationCommandOptionUser, Name: "usuario", Description: "Membro para consultar", Required: true},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "clear-warnings",
					Description: "Limpar avisos de um membro",
					Options: []*discordgo.ApplicationCommandOption{
						{Type: discordgo.ApplicationCommandOptionUser, Name: "usuario", Description: "Membro para limpar avisos", Required: true},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "list",
					Description: "Listar membros do servidor",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "search",
					Description: "Buscar membros por nome",
					Options: []*discordgo.ApplicationCommandOption{
						{Type: discordgo.ApplicationCommandOptionString, Name: "query", Description: "Termo de busca", Required: true},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "events",
					Description: "Ver eventos de membros recentes",
				},
			},
		},
		Handle: func(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
			data := i.ApplicationCommandData()
			sub := data.Options[0].Name

			switch sub {
			case "info":
				return handleMemberInfo(repo, data, i.GuildID)
			case "warn":
				return handleMemberWarn(repo, data, i.GuildID, i.ChannelID)
			case "warnings":
				return handleMemberWarnings(repo, data, i.GuildID)
			case "clear-warnings":
				return handleClearWarnings(repo, data, i.GuildID)
			case "list":
				return handleMemberList(repo, i.GuildID)
			case "search":
				return handleMemberSearch(repo, data)
			case "events":
				return handleMemberEvents(repo, i.GuildID)
			default:
				return "Subcomando desconhecido.", nil
			}
		},
	}
}

func handleMemberInfo(repo *sqlite.Repository, data discordgo.ApplicationCommandInteractionData, guildID string) (string, error) {
	user := data.Options[0].Options[0].UserValue(nil)
	if user == nil {
		return "Usuario nao encontrado.", nil
	}

	profile, err := repo.GetMemberProfile(user.ID)
	if err != nil {
		return fmt.Sprintf("Perfil nao encontrado para <@%s>. O usuario pode nao ter mensagens registradas.", user.ID), nil
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("**Perfil: %s**\n", profile.Name))
	b.WriteString(fmt.Sprintf("- ID: `%s`\n", profile.DiscordID))
	if len(profile.Roles) > 0 {
		b.WriteString(fmt.Sprintf("- Cargos: %s\n", strings.Join(profile.Roles, ", ")))
	}
	if profile.Age > 0 {
		b.WriteString(fmt.Sprintf("- Idade: %d\n", profile.Age))
	}
	if profile.Interests != "" {
		b.WriteString(fmt.Sprintf("- Interesses: %s\n", profile.Interests))
	}
	if profile.Notes != "" {
		b.WriteString(fmt.Sprintf("- Notas: %s\n", profile.Notes))
	}

	count, _ := repo.GetWarningCount(guildID, user.ID)
	b.WriteString(fmt.Sprintf("- Avisos: %d", count))

	return b.String(), nil
}

func handleMemberWarn(repo *sqlite.Repository, data discordgo.ApplicationCommandInteractionData, guildID, channelID string) (string, error) {
	user := data.Options[0].Options[0].UserValue(nil)
	if user == nil {
		return "Usuario nao encontrado.", nil
	}

	motivo := ""
	for _, opt := range data.Options[0].Options {
		if opt.Name == "motivo" {
			motivo = opt.StringValue()
		}
	}

	gc, err := repo.GetGuildConfig(guildID)
	if err != nil {
		gc = &sqlite.GuildConfig{MaxWarnings: 3}
	}

	if err := repo.AddWarning(guildID, user.ID, channelID, motivo, "", 1); err != nil {
		return "", err
	}

	count, _ := repo.GetWarningCount(guildID, user.ID)

	_ = repo.LogAuditEvent(guildID, "moderation", "warn", "", "Command", user.ID, user.Username,
		sqlite.AuditDetailsJSON(map[string]interface{}{
			"reason":         motivo,
			"warning_count":  count,
			"max_warnings":   gc.MaxWarnings,
		}), channelID)

	return fmt.Sprintf("**Aviso registrado**\n- Usuario: <@%s>\n- Motivo: %s\n- Avisos: %d/%d", user.ID, motivo, count, gc.MaxWarnings), nil
}

func handleMemberWarnings(repo *sqlite.Repository, data discordgo.ApplicationCommandInteractionData, guildID string) (string, error) {
	user := data.Options[0].Options[0].UserValue(nil)
	if user == nil {
		return "Usuario nao encontrado.", nil
	}

	warnings, err := repo.GetWarningsByUser(guildID, user.ID)
	if err != nil || len(warnings) == 0 {
		return fmt.Sprintf("<@%s> nao possui avisos.", user.ID), nil
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("**Avisos de %s** (%d total)\n", user.Username, len(warnings)))
	for i, w := range warnings {
		b.WriteString(fmt.Sprintf("%d. %s (gravidade: %d) - %s\n", i+1, w.Reason, w.Severity, w.CreatedAt))
	}

	return b.String(), nil
}

func handleClearWarnings(repo *sqlite.Repository, data discordgo.ApplicationCommandInteractionData, guildID string) (string, error) {
	user := data.Options[0].Options[0].UserValue(nil)
	if user == nil {
		return "Usuario nao encontrado.", nil
	}

	_, err := repo.GetDB().Exec("DELETE FROM guild_warnings WHERE guild_id = ? AND user_id = ?", guildID, user.ID)
	if err != nil {
		return "", err
	}

	_ = repo.LogAuditEvent(guildID, "moderation", "clear_warnings", "", "Command", user.ID, user.Username, "", "")

	return fmt.Sprintf("Avisos de <@%s> foram limpos.", user.ID), nil
}

func handleMemberList(repo *sqlite.Repository, guildID string) (string, error) {
	profiles, err := repo.ListMemberProfiles(20, 0)
	if err != nil || len(profiles) == 0 {
		return "Nenhum membro encontrado.", nil
	}

	var b strings.Builder
	b.WriteString("**Membros Registrados**\n")
	for _, p := range profiles {
		count, _ := repo.GetWarningCount(guildID, p.DiscordID)
		b.WriteString(fmt.Sprintf("- %s (`%s`) - %d avisos\n", p.Name, p.DiscordID, count))
	}

	total, _ := repo.CountMemberProfiles()
	if total > 20 {
		b.WriteString(fmt.Sprintf("\n... e mais %d membros", total-20))
	}

	return b.String(), nil
}

func handleMemberSearch(repo *sqlite.Repository, data discordgo.ApplicationCommandInteractionData) (string, error) {
	query := ""
	for _, opt := range data.Options[0].Options {
		if opt.Name == "query" {
			query = opt.StringValue()
		}
	}

	profiles, err := repo.SearchMemberProfiles(query)
	if err != nil || len(profiles) == 0 {
		return "Nenhum membro encontrado.", nil
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("**Resultados para '%s'**\n", query))
	for _, p := range profiles {
		b.WriteString(fmt.Sprintf("- %s (`%s`)\n", p.Name, p.DiscordID))
	}

	return b.String(), nil
}

func handleMemberEvents(repo *sqlite.Repository, guildID string) (string, error) {
	events, err := repo.GetMemberEvents(guildID, 10)
	if err != nil || len(events) == 0 {
		return "Nenhum evento recente.", nil
	}

	var b strings.Builder
	b.WriteString("**Eventos Recentes**\n")
	for _, e := range events {
		icon := "➡️"
		switch e.EventType {
		case "leave":
			icon = "⬅️"
		case "kick":
			icon = "👢"
		case "ban":
			icon = "🔨"
		}
		b.WriteString(fmt.Sprintf("%s %s - %s (%s)\n", icon, e.UserName, e.EventType, e.CreatedAt))
	}

	return b.String(), nil
}
