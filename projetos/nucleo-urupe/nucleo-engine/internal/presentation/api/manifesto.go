/*
 * Copyright (c) 2026 Frente Urupê
 * Licensed under the MIT License.
 */

package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"nucleo-engine/internal/data/sqlite"
)

func handleManifestoVersionsList(repo *sqlite.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		versions, err := repo.ListManifestoVersions()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if versions == nil {
			versions = []sqlite.ManifestoVersion{}
		}

		json.NewEncoder(w).Encode(versions)
	}
}

func handleManifestoActiveGet(repo *sqlite.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		manifesto, err := repo.GetActiveManifesto()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if manifesto == nil {
			http.Error(w, "No active manifesto version found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(manifesto)
	}
}

func handleManifestoVersionSave(repo *sqlite.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req sqlite.ManifestoVersion
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if req.Version == "" || req.Title == "" || req.Content == "" {
			http.Error(w, "Version, Title and Content are required", http.StatusBadRequest)
			return
		}

		if req.ID == "" {
			req.ID = "man-" + uuid.New().String()[:8]
		}
		if req.Author == "" {
			req.Author = "Coordenação Nacional"
		}

		if err := repo.CreateManifestoVersion(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"version": req,
		})
	}
}

func handleManifestoVersionActivate(repo *sqlite.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id parameter is required", http.StatusBadRequest)
			return
		}

		if err := repo.ActivateManifestoVersion(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}
