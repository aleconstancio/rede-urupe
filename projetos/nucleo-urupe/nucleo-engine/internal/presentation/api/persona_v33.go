/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	minotaur "github.com/aleconstancio/talos/v2/behavior/persona"

	"github.com/google/uuid"
)

func (s *Server) handleIdentities(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := s.repo.ListCoreIdentityProfiles()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"items": items})
	case http.MethodPost:
		var item minotaur.CoreIdentityProfile
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if strings.TrimSpace(item.ID) == "" {
			item.ID = uuid.NewString()
		}
		if err := s.repo.UpsertCoreIdentityProfile(item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		s.Broadcast("persona")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "item": item})
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handleIdentityByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/identities/"))
	if id == "" {
		http.Error(w, "identity id required", 400)
		return
	}

	switch r.Method {
	case http.MethodGet:
		item, err := s.repo.GetCoreIdentityProfile(id)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(item)
	case http.MethodPost:
		var item minotaur.CoreIdentityProfile
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		item.ID = id
		if err := s.repo.UpsertCoreIdentityProfile(item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		s.Broadcast("persona")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "item": item})
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handlePersonas(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		identityID := strings.TrimSpace(r.URL.Query().Get("identity_id"))
		items, err := s.repo.ListPersonaOverlays(identityID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"items": items})
	case http.MethodPost:
		var item minotaur.PersonaOverlay
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if strings.TrimSpace(item.ID) == "" {
			item.ID = uuid.NewString()
		}
		if err := s.repo.UpsertPersonaOverlay(item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		s.Broadcast("persona")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "item": item})
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handlePersonaByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/personas/"))
	if id == "" {
		http.Error(w, "persona id required", 400)
		return
	}

	switch r.Method {
	case http.MethodGet:
		item, err := s.repo.GetPersonaOverlay(id)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		json.NewEncoder(w).Encode(item)
	case http.MethodPost:
		var item minotaur.PersonaOverlay
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		item.ID = id
		if err := s.repo.UpsertPersonaOverlay(item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		s.Broadcast("persona")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "item": item})
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handlePersonaPolicy(w http.ResponseWriter, r *http.Request) {
	channelID := s.personaChannelID(r)
	switch r.Method {
	case http.MethodGet:
		item, err := s.repo.GetPersonaPolicy(channelID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(item)
	case http.MethodPost:
		var item minotaur.PersonaPolicy
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if strings.TrimSpace(item.ChannelID) == "" {
			item.ChannelID = channelID
		}
		if err := s.repo.UpsertPersonaPolicy(item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		s.Broadcast("persona")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "item": item})
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handlePersonaAdaptations(w http.ResponseWriter, r *http.Request) {
	channelID := s.personaChannelID(r)
	switch r.Method {
	case http.MethodGet:
		items, err := s.repo.ListAdaptivePersonaMemories(channelID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"items": items})
	case http.MethodPost:
		var item minotaur.AdaptivePersonaMemory
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if strings.TrimSpace(item.ChannelID) == "" {
			item.ChannelID = channelID
		}
		if err := s.repo.UpsertAdaptivePersonaMemory(item); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		s.Broadcast("persona")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "item": item})
	default:
		http.Error(w, "Method not allowed", 405)
	}
}

func (s *Server) handlePersonaAdaptationsReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var body struct {
		ChannelID  string `json:"channel_id"`
		IdentityID string `json:"identity_id"`
		PersonaID  string `json:"persona_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if strings.TrimSpace(body.ChannelID) == "" {
		body.ChannelID = s.personaChannelID(r)
	}
	if err := s.repo.ResetAdaptivePersonaMemories(body.ChannelID, body.IdentityID, body.PersonaID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s.Broadcast("persona")
	json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func (s *Server) handlePersonaProposals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", 405)
		return
	}
	channelID := s.personaChannelID(r)
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	items, err := s.repo.ListPersonaDeltaProposals(channelID, status)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"items": items})
}

func (s *Server) handlePersonaProposalAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	path := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/persona-proposals/"))
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		http.Error(w, "invalid proposal path", 400)
		return
	}
	id := strings.TrimSpace(parts[0])
	action := strings.TrimSpace(parts[1])
	if id == "" {
		http.Error(w, "proposal id required", 400)
		return
	}

	if action == "apply" {
		if err := s.repo.ApplyPersonaProposal(id); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		s.Broadcast("persona")
		json.NewEncoder(w).Encode(map[string]any{"ok": true, "status": "applied"})
		return
	}

	status, err := proposalActionToStatus(action)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := s.repo.UpdatePersonaDeltaProposalStatus(id, status); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	s.Broadcast("persona")
	json.NewEncoder(w).Encode(map[string]any{"ok": true, "status": status})
}

func (s *Server) personaChannelID(r *http.Request) string {
	channelID := strings.TrimSpace(r.URL.Query().Get("channel_id"))
	if channelID == "" {
		channelID = s.cfg.TargetChannelID
	}
	return channelID
}

func proposalActionToStatus(action string) (string, error) {
	switch action {
	case "approve":
		return "approved", nil
	case "reject":
		return "rejected", nil
	case "apply":
		return "applied", nil
	default:
		return "", errors.New("unsupported proposal action")
	}
}
