package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"nucleo-engine/internal/data/sqlite"
)

func (s *Server) handleShowcaseProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		projects, err := s.repo.ListShowcaseProjects()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if projects == nil {
			projects = []sqlite.ShowcaseProject{}
		}
		json.NewEncoder(w).Encode(map[string]any{"projects": projects})

	case http.MethodPost:
		var p sqlite.ShowcaseProject
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		p.Slug = strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))
		if err := s.repo.CreateShowcaseProject(&p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.Broadcast("projects")
		json.NewEncoder(w).Encode(p)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleShowcaseProjectBySlug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	slug := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if slug == "" {
		http.Error(w, "slug required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		p, err := s.repo.GetShowcaseProject(slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if p == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(p)

	case http.MethodPut:
		var p sqlite.ShowcaseProject
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		p.Slug = slug
		if err := s.repo.UpdateShowcaseProject(&p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.Broadcast("projects")
		json.NewEncoder(w).Encode(p)

	case http.MethodDelete:
		if err := s.repo.DeleteShowcaseProject(slug); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.Broadcast("projects")
		json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
