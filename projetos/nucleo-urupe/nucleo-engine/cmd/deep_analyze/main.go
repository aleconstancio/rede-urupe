/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"nucleo-engine/internal/config"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type StatPair struct {
	Key   string
	Value int
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open("sqlite3", cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	fmt.Println("# Discord Agent Deep Intelligence Report")
	fmt.Printf("Analysis Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("---")

	// 1. Executive Summary
	var total, botMsgs, userMsgs int
	var earliest, latest string
	db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&total)
	db.QueryRow("SELECT COUNT(*) FROM messages WHERE is_bot = 1").Scan(&botMsgs)
	db.QueryRow("SELECT COUNT(*) FROM messages WHERE is_bot = 0").Scan(&userMsgs)
	db.QueryRow("SELECT MIN(timestamp), MAX(timestamp) FROM messages").Scan(&earliest, &latest)

	fmt.Println("## 1. Executive Summary")
	fmt.Printf("- **Total Volume**: %d messages\n", total)
	fmt.Printf("- **Human Participation**: %d messages\n", userMsgs)
	fmt.Printf("- **Agent Activity**: %d messages (%.1f%% reactive rate)\n", botMsgs, float64(botMsgs)/float64(userMsgs)*100)
	fmt.Printf("- **Data Range**: From %s to %s\n\n", earliest[:16], latest[:16])

	// 2. Temporal Activity (Heatmap)
	fmt.Println("## 2. Temporal Engagement Trends")
	fmt.Println("### Hourly Distribution (All time)")
	hourly := make(map[int]int)
	rows, _ := db.Query("SELECT strftime('%H', timestamp) as hour FROM messages WHERE is_bot = 0")
	for rows.Next() {
		var hour int
		rows.Scan(&hour)
		hourly[hour]++
	}
	rows.Close()

	maxH := 0
	for _, v := range hourly {
		if v > maxH {
			maxH = v
		}
	}

	for h := 0; h < 24; h++ {
		count := hourly[h]
		barLen := 0
		if maxH > 0 {
			barLen = (count * 30) / maxH
		}
		fmt.Printf("%02d:00 | %-30s (%d)\n", h, strings.Repeat("█", barLen), count)
	}
	fmt.Println()

	// 3. User Power Ranking
	fmt.Println("## 3. Top Community Drivers")
	rows, _ = db.Query(`
		SELECT author, COUNT(*) as count, COUNT(DISTINCT channel_id) as channels
		FROM messages 
		WHERE is_bot = 0 
		GROUP BY author 
		ORDER BY count DESC 
		LIMIT 15
	`)
	fmt.Println("| User | Message Count | Diversity (Channels) |")
	fmt.Println("| :--- | :--- | :--- |")
	for rows.Next() {
		var author string
		var count, channels int
		rows.Scan(&author, &count, &channels)
		fmt.Printf("| %s | %d | %d |\n", author, count, channels)
	}
	rows.Close()
	fmt.Println()

	// 4. Content Network (Mentions)
	fmt.Println("## 4. Interaction Network (Top Mentions)")
	mentionCounts := make(map[string]int)
	rows, _ = db.Query("SELECT content FROM messages WHERE is_bot = 0 AND content LIKE '%<@%'")
	for rows.Next() {
		var content string
		rows.Scan(&content)
		words := strings.Fields(content)
		for _, w := range words {
			if strings.HasPrefix(w, "<@") && strings.HasSuffix(w, ">") {
				mentionCounts[w]++
			}
		}
	}
	rows.Close()

	var mentions []StatPair
	for k, v := range mentionCounts {
		mentions = append(mentions, StatPair{k, v})
	}
	sort.Slice(mentions, func(i, j int) bool { return mentions[i].Value > mentions[j].Value })

	for i := 0; i < 10 && i < len(mentions); i++ {
		fmt.Printf("- %s: mentioned %d times\n", mentions[i].Key, mentions[i].Value)
	}
	fmt.Println()

	// 5. Emerging Keywords (Daily Trend)
	fmt.Println("## 5. Daily Topic Evolution")
	stopWords := map[string]bool{"a": true, "o": true, "e": true, "que": true, "do": true, "da": true, "de": true, "um": true, "uma": true, "com": true, "em": true, "no": true, "na": true, "para": true, "por": true, "é": true, "não": true, "mais": true, "como": true, "você": true, "isso": true}
	
	dailyWords := make(map[string]map[string]int)
	rows, _ = db.Query("SELECT strftime('%Y-%m-%d', timestamp) as day, content FROM messages WHERE is_bot = 0")
	for rows.Next() {
		var day, content string
		rows.Scan(&day, &content)
		if dailyWords[day] == nil {
			dailyWords[day] = make(map[string]int)
		}
		words := strings.Fields(strings.ToLower(content))
		for _, w := range words {
			w = strings.Trim(w, ".,!?;:()\"")
			if len(w) > 4 && !stopWords[w] && !strings.HasPrefix(w, "<@") {
				dailyWords[day][w]++
			}
		}
	}
	rows.Close()

	var days []string
	for d := range dailyWords {
		days = append(days, d)
	}
	sort.Strings(days)

	for _, d := range days {
		var words []StatPair
		for w, c := range dailyWords[d] {
			words = append(words, StatPair{w, c})
		}
		sort.Slice(words, func(i, j int) bool { return words[i].Value > words[j].Value })
		
		top := []string{}
		for i := 0; i < 5 && i < len(words); i++ {
			top = append(top, fmt.Sprintf("%s (%d)", words[i].Key, words[i].Value))
		}
		fmt.Printf("- **%s**: %s\n", d, strings.Join(top, ", "))
	}
	fmt.Println()

	// 6. Channel Density
	fmt.Println("## 6. Channel Activity Distribution")
	rows, _ = db.Query("SELECT channel_id, COUNT(*) as count FROM messages GROUP BY channel_id ORDER BY count DESC")
	for rows.Next() {
		var cid string
		var count int
		rows.Scan(&cid, &count)
		fmt.Printf("- Channel %s: %d messages (%.1f%% total)\n", cid, count, float64(count)/float64(total)*100)
	}
	rows.Close()
}
