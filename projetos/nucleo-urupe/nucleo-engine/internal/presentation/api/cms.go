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

func handleCMSArticlesList(repo *sqlite.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		pubOnly := r.URL.Query().Get("public") == "true"

		articles, err := repo.ListArticles(pubOnly)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if articles == nil {
			articles = []sqlite.Article{}
		}

		json.NewEncoder(w).Encode(articles)
	}
}

func handleCMSArticleGet(repo *sqlite.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		slug := r.URL.Query().Get("slug")
		if slug == "" {
			http.Error(w, "slug parameter is required", http.StatusBadRequest)
			return
		}

		article, err := repo.GetArticleBySlug(slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if article == nil {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(article)
	}
}

func handleCMSArticleSave(repo *sqlite.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req sqlite.Article
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.Content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		if req.ID == "" {
			req.ID = "art-" + uuid.New().String()[:8]
		}
		if req.Author == "" {
			req.Author = "Frente Urupê"
		}
		if req.Category == "" {
			req.Category = "Notícia"
		}

		if err := repo.UpsertArticle(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"article": req,
		})
	}
}

func handleCMSArticleDelete(repo *sqlite.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "id parameter is required", http.StatusBadRequest)
			return
		}

		if err := repo.DeleteArticle(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}
