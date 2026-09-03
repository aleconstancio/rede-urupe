/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */

package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"nucleo-engine/internal/data/sqlite"
)

func (s *Server) handleForumTemplates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		guildID := r.URL.Query().Get("guild_id")
		id := r.URL.Query().Get("id")
		if id != "" {
			t, err := s.repo.GetForumTemplate(id)
			if err != nil {
				http.Error(w, err.Error(), 404)
				return
			}
			json.NewEncoder(w).Encode(t)
			return
		}
		templates, err := s.repo.ListForumTemplates(guildID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(templates)

	case http.MethodPost:
		var t sqlite.ForumTemplate
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		if t.ID == "" {
			t.ID = uuid.New().String()
		}
		if err := s.repo.UpsertForumTemplate(t); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(t)

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "missing id", 400)
			return
		}
		if err := s.repo.DeleteForumTemplate(id); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(204)

	default:
		http.Error(w, "method not allowed", 405)
	}
}

func (s *Server) handleForumPublish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing template id", 400)
		return
	}

	t, err := s.repo.GetForumTemplate(id)
	if err != nil {
		http.Error(w, "template not found", 404)
		return
	}

	var vars map[string]string
	if r.ContentLength > 0 {
		json.NewDecoder(r.Body).Decode(&vars)
	}

	if s.forumWorker == nil {
		json.NewEncoder(w).Encode(map[string]any{
			"status": "failed",
			"error":  "forum worker not initialized",
		})
		return
	}
	post, err := s.forumWorker.PublishFromTemplate(r.Context(), *t, vars)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"status": "failed",
			"error":  err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (s *Server) handleForumPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}

	guildID := r.URL.Query().Get("guild_id")
	posts, err := s.repo.ListForumPosts(guildID, 50)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(posts)
}

// extractIDFromPath parses /api/forum/templates/some-uuid from the path
func extractIDFromPath(prefix, path string) string {
	trimmed := strings.TrimPrefix(path, prefix)
	trimmed = strings.Trim(trimmed, "/")
	return trimmed
}
