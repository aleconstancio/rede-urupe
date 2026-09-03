/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"nucleo-engine/internal/data/sqlite"
)

func (s *Server) handleAdminMembers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		limit := 50
		offset := 0
		if q := r.URL.Query().Get("limit"); q != "" {
			if v, err := parseInt(q); err == nil {
				limit = v
			}
		}
		if q := r.URL.Query().Get("offset"); q != "" {
			if v, err := parseInt(q); err == nil {
				offset = v
			}
		}

		profiles, err := s.repo.ListMemberProfiles(limit, offset)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		total, _ := s.repo.CountMemberProfiles()

		json.NewEncoder(w).Encode(map[string]interface{}{
			"members": profiles,
			"total":   total,
			"limit":   limit,
			"offset":  offset,
		})
	case http.MethodPost:
		var p sqlite.MemberProfile
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if err := s.repo.UpsertMemberProfile(&p); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handleAdminMemberByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/member/")
	if id == "" {
		http.Error(w, "Missing member ID", 400)
		return
	}

	switch r.Method {
	case http.MethodGet:
		profile, err := s.repo.GetMemberProfile(id)
		if err != nil {
			http.Error(w, "Member not found", 404)
			return
		}
		json.NewEncoder(w).Encode(profile)
	case http.MethodPost:
		var p sqlite.MemberProfile
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		p.DiscordID = id
		if err := s.repo.UpsertMemberProfile(&p); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handleAdminAudit(w http.ResponseWriter, r *http.Request) {
	guildID := s.cfg.TargetGuildID
	limit := 50
	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := parseInt(q); err == nil {
			limit = v
		}
	}

	events, err := s.repo.GetAuditLog(guildID, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
		"total":  len(events),
	})
}

func (s *Server) handleAdminModLog(w http.ResponseWriter, r *http.Request) {
	guildID := s.cfg.TargetGuildID
	limit := 50
	if q := r.URL.Query().Get("limit"); q != "" {
		if v, err := parseInt(q); err == nil {
			limit = v
		}
	}

	warnings, err := s.repo.GetAuditLogByAction(guildID, "moderation", limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": warnings,
		"total":  len(warnings),
	})
}

func (s *Server) handleAdminWelcome(w http.ResponseWriter, r *http.Request) {
	guildID := s.cfg.TargetGuildID

	switch r.Method {
	case http.MethodGet:
		var wc struct {
			Enabled         bool   `json:"enabled"`
			ChannelID       string `json:"channel_id"`
			WelcomeMessage  string `json:"welcome_message"`
			GoodbyeMessage  string `json:"goodbye_message"`
		}
		err := s.repo.GetDB().QueryRow(`
			SELECT enabled, channel_id, welcome_message, goodbye_message
			FROM welcome_config WHERE guild_id = ?
		`, guildID).Scan(&wc.Enabled, &wc.ChannelID, &wc.WelcomeMessage, &wc.GoodbyeMessage)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"enabled": false,
			})
			return
		}
		json.NewEncoder(w).Encode(wc)
	case http.MethodPost:
		var req struct {
			ChannelID      string `json:"channel_id"`
			WelcomeMessage string `json:"welcome_message"`
			GoodbyeMessage string `json:"goodbye_message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		_, err := s.repo.GetDB().Exec(`
			INSERT INTO welcome_config (guild_id, enabled, channel_id, welcome_message, goodbye_message, updated_at)
			VALUES (?, 1, ?, ?, ?, datetime('now'))
			ON CONFLICT(guild_id) DO UPDATE SET
				enabled = 1,
				channel_id = excluded.channel_id,
				welcome_message = excluded.welcome_message,
				goodbye_message = excluded.goodbye_message,
				updated_at = excluded.updated_at
		`, guildID, req.ChannelID, req.WelcomeMessage, req.GoodbyeMessage)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func parseInt(s string) (int, error) {
	var v int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		v = v*10 + int(c-'0')
	}
	return v, nil
}
