/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */

package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"nucleo-engine/internal/data/sqlite"
)

func StatusCommand(repo *sqlite.Repository) *Cmd {
	return &Cmd{
		Def: &discordgo.ApplicationCommand{
			Name:        "urupe",
			Description: "Status e informacoes do Núcleo Urupê",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "status",
					Description: "Status atual do Núcleo Urupê e da mascote Micélia 🍄",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "info",
					Description: "Informacoes do servidor e canal ativo",
				},
			},
		},
		Handle: func(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
			data := i.ApplicationCommandData()
			sub := data.Options[0].Name

			switch sub {
			case "status":
				activePersona := "Micélia 🍄"
				studio, err := repo.GetPersonaStudioState(i.GuildID)
				if err == nil && studio.ActiveIdentity.Name != "" {
					activePersona = studio.ActiveIdentity.Name
				}
				return fmt.Sprintf("**Núcleo Urupê** ativo\n- Persona: %s\n- Módulo: Cognitivo (Pulse & Gateway)\n- Memória: SQLite FTS5\n- Servidor: `%s`", activePersona, i.GuildID), nil

			case "info":
				gc, err := repo.GetGuildConfig(i.GuildID)
				if err != nil {
					// Default info if no guild config
					return fmt.Sprintf("Servidor: `%s`\nCanal padrao: nenhum configurado\nUse `/urupe-config` para configurar.", i.GuildID), nil
				}
				return fmt.Sprintf("**Config do Servidor**\n- Canal padrao: <#%s>\n- Prefixo: `%s`\n- Mod Log: <#%s>\n- Moderacao: %s",
					gc.DefaultChannelID, gc.Prefix, gc.ModLogChannelID, boolStr(gc.ModEnabled)), nil

			default:
				return "Subcomando desconhecido.", nil
			}
		},
	}
}

func ConfigCommand(repo *sqlite.Repository) *Cmd {
	return &Cmd{
		Def: &discordgo.ApplicationCommand{
			Name:        "urupe-config",
			Description: "Gerenciar configuracao do Núcleo Urupê no servidor",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "show",
					Description: "Mostrar configuracao atual",
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set-channel",
					Description: "Definir canal padrao",
					Options: []*discordgo.ApplicationCommandOption{
						{Type: discordgo.ApplicationCommandOptionChannel, Name: "canal", Description: "Canal do Discord", Required: true},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "set-modlog",
					Description: "Definir canal de moderacao",
					Options: []*discordgo.ApplicationCommandOption{
						{Type: discordgo.ApplicationCommandOptionChannel, Name: "canal", Description: "Canal de mod log", Required: true},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "toggle-mod",
					Description: "Ativar/desativar moderacao automatica",
				},
			},
		},
		Handle: func(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error) {
			data := i.ApplicationCommandData()
			sub := data.Options[0].Name

			gc, err := repo.GetGuildConfig(i.GuildID)
			if err != nil {
				gc = &sqlite.GuildConfig{
					GuildID:     i.GuildID,
					Prefix:      "/",
					ModEnabled:  false,
					MaxWarnings: 3,
				}
			}

			switch sub {
			case "show":
				return fmt.Sprintf("**Configuracao**\n- Canal: <#%s>\n- Mod Log: <#%s>\n- Moderacao: %s\n- Max Avisos: %d",
					gc.DefaultChannelID, gc.ModLogChannelID, boolStr(gc.ModEnabled), gc.MaxWarnings), nil

			case "set-channel":
				chID := data.Options[0].Options[0].ChannelValue(s).ID
				gc.DefaultChannelID = chID
				if err := repo.UpsertGuildConfig(*gc); err != nil {
					return "", err
				}
				return fmt.Sprintf("Canal padrao definido para <#%s>", chID), nil

			case "set-modlog":
				chID := data.Options[0].Options[0].ChannelValue(s).ID
				gc.ModLogChannelID = chID
				if err := repo.UpsertGuildConfig(*gc); err != nil {
					return "", err
				}
				return fmt.Sprintf("Canal de mod log definido para <#%s>", chID), nil

			case "toggle-mod":
				gc.ModEnabled = !gc.ModEnabled
				if err := repo.UpsertGuildConfig(*gc); err != nil {
					return "", err
				}
				return fmt.Sprintf("Moderacao automatica: %s", boolStr(gc.ModEnabled)), nil

			default:
				return "Subcomando desconhecido.", nil
			}
		},
	}
}

func boolStr(b bool) string {
	if b {
		return "ativada"
	}
	return "desativada"
}
