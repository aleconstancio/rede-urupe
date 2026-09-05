/*
 * Copyright (c) 2026 Frente Urupê
 * Licensed under the MIT License.
 */

package api

import (
	"encoding/json"
	"net/http"

	"nucleo-engine/internal/data/sqlite"
)

func handleMiceliaMemoriesList(repo *sqlite.Repository, defaultChannelID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		channelID := r.URL.Query().Get("channel_id")
		if channelID == "" {
			channelID = defaultChannelID
		}

		annotations, err := repo.ListMessageAnnotations(channelID, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if annotations == nil {
			annotations = []sqlite.MessageAnnotation{}
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"channel_id":  channelID,
			"annotations": annotations,
		})
	}
}

func handleMiceliaMemoryDelete(repo *sqlite.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		messageID := r.URL.Query().Get("message_id")
		if messageID == "" {
			http.Error(w, "message_id parameter is required", http.StatusBadRequest)
			return
		}

		if err := repo.DeleteMessageAnnotation(messageID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	}
}
