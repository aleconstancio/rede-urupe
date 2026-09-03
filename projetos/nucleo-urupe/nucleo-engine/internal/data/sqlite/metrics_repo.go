/*
 * Copyright (c) 2026 Talos V2 Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"encoding/json"
	"time"
	"nucleo-engine/internal/pkg/timeutil"
)

// DashboardStats contains all the aggregations needed for the IntelligencePanel.
type DashboardStats struct {
	Hourly     []StatPair  `json:"hourly"`
	Categories []StatPair  `json:"categories"`
	Authors    []StatPair  `json:"authors"`
	Heatmap    []HeatPoint `json:"heatmap"`
	Keywords   []StatPair  `json:"keywords"`
	System     struct {
		APICalls      int `json:"api_calls"`
		MessagesSent  int `json:"messages_sent"`
		ScrapingRuns  int `json:"scraping_runs"`
		AvgLatency    int `json:"avg_latency"`
		UptimeMinutes int `json:"uptime_minutes"`
	} `json:"system"`
}

type StatPair struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type HeatPoint struct {
	Day   int `json:"day"`
	Hour  int `json:"hour"`
	Value int `json:"value"`
}

func (r *Repository) GetDashboardStats(channelID string) (DashboardStats, error) {
	var stats DashboardStats

	// 1. Hourly (Last 24h)
	dayAgo := timeutil.Now().Add(-24 * time.Hour)
	rows, err := r.db.Query(`
		SELECT strftime('%H:00', timestamp) as hr, COUNT(*) 
		FROM messages 
		WHERE channel_id = ? AND timestamp > ?
		GROUP BY hr ORDER BY hr ASC`, channelID, dayAgo)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p StatPair
			if err := rows.Scan(&p.Key, &p.Value); err == nil {
				stats.Hourly = append(stats.Hourly, p)
			}
		}
	}

	// 2. Categories
	rows, err = r.db.Query(`
		SELECT category, COUNT(*) 
		FROM messages 
		WHERE channel_id = ? 
		GROUP BY category ORDER BY COUNT(*) DESC LIMIT 10`, channelID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p StatPair
			if err := rows.Scan(&p.Key, &p.Value); err == nil {
				stats.Categories = append(stats.Categories, p)
			}
		}
	}

	// 3. Authors
	rows, err = r.db.Query(`
		SELECT author, COUNT(*) 
		FROM messages 
		WHERE channel_id = ? 
		GROUP BY author ORDER BY COUNT(*) DESC LIMIT 10`, channelID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p StatPair
			if err := rows.Scan(&p.Key, &p.Value); err == nil {
				stats.Authors = append(stats.Authors, p)
			}
		}
	}

	// 4. Heatmap
	rows, err = r.db.Query(`
		SELECT CAST(strftime('%w', timestamp) AS INT) as day, 
		       CAST(strftime('%H', timestamp) AS INT) as hr, 
		       COUNT(*)
		FROM messages 
		WHERE channel_id = ?
		GROUP BY day, hr`, channelID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var p HeatPoint
			if err := rows.Scan(&p.Day, &p.Hour, &p.Value); err == nil {
				stats.Heatmap = append(stats.Heatmap, p)
			}
		}
	}

	// 5. System Metrics
	dayAgo = timeutil.Now().Add(-24 * time.Hour)
	r.db.QueryRow(`SELECT COUNT(*) FROM system_metrics WHERE metric_type = 'api_call' AND timestamp > ?`, dayAgo).Scan(&stats.System.APICalls)
	r.db.QueryRow(`SELECT COUNT(*) FROM messages WHERE is_bot = 1 AND timestamp > ?`, dayAgo).Scan(&stats.System.MessagesSent)
	r.db.QueryRow(`SELECT COUNT(*) FROM system_metrics WHERE metric_type = 'scraping_run' AND timestamp > ?`, dayAgo).Scan(&stats.System.ScrapingRuns)
	r.db.QueryRow(`SELECT CAST(AVG(latency_ms) AS INTEGER) FROM system_metrics WHERE timestamp > ?`, dayAgo).Scan(&stats.System.AvgLatency)
	
	// Simulate uptime (based on first message)
	var firstMsgAt time.Time
	if err := r.db.QueryRow(`SELECT MIN(timestamp) FROM messages`).Scan(&firstMsgAt); err == nil {
		stats.System.UptimeMinutes = int(time.Since(firstMsgAt).Minutes())
	}

	return stats, nil
}

func (r *Repository) LogMetric(mType, subType string, value float64, latency int, metadata map[string]any) error {
	metaJSON := "{}"
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			metaJSON = string(b)
		}
	}

	_, err := r.db.Exec(`
		INSERT INTO system_metrics (metric_type, sub_type, value, latency_ms, metadata_json)
		VALUES (?, ?, ?, ?, ?)
	`, mType, subType, value, latency, metaJSON)
	return err
}

type SystemSummary struct {
	APICalls     int
	MessagesSent int
	ScrapingRuns int
	AvgLatency   int
}

func (r *Repository) GetSystemSummary(channelID string) (SystemSummary, error) {
	var s SystemSummary

	// Optimized consolidated query
	dayAgo := timeutil.Now().Add(-24 * time.Hour)
	err := r.db.QueryRow(`
		SELECT 
			(SELECT COUNT(*) FROM system_metrics WHERE metric_type = 'api_call' AND timestamp > ?),
			(SELECT COUNT(*) FROM messages WHERE is_bot = 1 AND channel_id = ? AND timestamp > ?),
			(SELECT COUNT(*) FROM system_metrics WHERE metric_type = 'scraping_run' AND timestamp > ?),
			(SELECT COALESCE(CAST(AVG(latency_ms) AS INTEGER), 0) FROM system_metrics WHERE metric_type = 'api_call' AND timestamp > ?)
	`, dayAgo, channelID, dayAgo, dayAgo, dayAgo).Scan(&s.APICalls, &s.MessagesSent, &s.ScrapingRuns, &s.AvgLatency)

	return s, err
}
