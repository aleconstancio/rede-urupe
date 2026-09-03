/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleAnalyticsSentiment(w http.ResponseWriter, r *http.Request) {
	channelID := s.cfg.TargetChannelID
	days := 30
	if q := r.URL.Query().Get("days"); q != "" {
		if v, err := parseInt(q); err == nil && v > 0 {
			days = v
		}
	}

	data, err := s.repo.GetSentimentTrend(channelID, days)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"channel_id": channelID,
		"days":       days,
		"sentiment":  data,
	})
}

func (s *Server) handleAnalyticsTokens(w http.ResponseWriter, r *http.Request) {
	days := 30
	if q := r.URL.Query().Get("days"); q != "" {
		if v, err := parseInt(q); err == nil && v > 0 {
			days = v
		}
	}

	byModel, _ := s.repo.GetTokenUsageByModel(days)
	byReason, _ := s.repo.GetTokenUsageByReason(days)
	trend, _ := s.repo.GetTokenUsageTrend(days)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"days":     days,
		"by_model": byModel,
		"by_reason": byReason,
		"trend":    trend,
	})
}

func (s *Server) handleAnalyticsGrowth(w http.ResponseWriter, r *http.Request) {
	guildID := s.cfg.TargetGuildID
	days := 30
	if q := r.URL.Query().Get("days"); q != "" {
		if v, err := parseInt(q); err == nil && v > 0 {
			days = v
		}
	}

	data, err := s.repo.GetMemberGrowthDaily(guildID, days)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"guild_id": guildID,
		"days":     days,
		"growth":   data,
	})
}

func (s *Server) handleAnalyticsChannels(w http.ResponseWriter, r *http.Request) {
	health, err := s.repo.GetChannelHealth(s.cfg.TargetChannelID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"channel": health,
	})
}

func (s *Server) handleAnalyticsEngagement(w http.ResponseWriter, r *http.Request) {
	channelID := s.cfg.TargetChannelID
	days := 7
	if q := r.URL.Query().Get("days"); q != "" {
		if v, err := parseInt(q); err == nil && v > 0 {
			days = v
		}
	}

	score, err := s.repo.GetEngagementScore(channelID, days)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"channel_id": channelID,
		"days":       days,
		"score":      score,
	})
}
