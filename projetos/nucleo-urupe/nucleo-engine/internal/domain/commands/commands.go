/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */

// Package commands implements a command router and Discord slash command definitions.
package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Cmd struct {
	Def    *discordgo.ApplicationCommand
	Handle func(s *discordgo.Session, i *discordgo.InteractionCreate) (string, error)
}

type Router struct {
	commands map[string]*Cmd
}

func NewRouter() *Router {
	return &Router{commands: make(map[string]*Cmd)}
}

func (r *Router) Register(cmd *Cmd) {
	r.commands[cmd.Def.Name] = cmd
}

func (r *Router) AllCommands() []*discordgo.ApplicationCommand {
	var defs []*discordgo.ApplicationCommand
	for _, cmd := range r.commands {
		defs = append(defs, cmd.Def)
	}
	return defs
}

func (r *Router) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	cmd, ok := r.commands[data.Name]
	if !ok {
		return
	}

	msg, err := cmd.Handle(s, i)
	if err != nil {
		msg = "Erro interno: " + err.Error()
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
