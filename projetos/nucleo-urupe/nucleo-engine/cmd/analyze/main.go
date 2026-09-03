/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"nucleo-engine/internal/config"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

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

	fmt.Println("# Discord Agent Data Analysis")
	fmt.Printf("Generated on: %s\n\n", time.Now().Format(time.RFC3339))

	// 1. General Stats
	var total, botMsgs, userMsgs int
	db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&total)
	db.QueryRow("SELECT COUNT(*) FROM messages WHERE is_bot = 1").Scan(&botMsgs)
	db.QueryRow("SELECT COUNT(*) FROM messages WHERE is_bot = 0").Scan(&userMsgs)

	fmt.Println("## 1. Volume & Engagement")
	fmt.Printf("- Total Messages: %d\n", total)
	fmt.Printf("- User Messages: %d\n", userMsgs)
	fmt.Printf("- Bot Responses: %d (%.2f%% engagement)\n\n", botMsgs, float64(botMsgs)/float64(total)*100)

	// 2. Persona Check
	fmt.Println("## 2. Active Persona Status")
	var name, model string
	var temp float64
	err = db.QueryRow("SELECT name, model, temperature FROM agent_profiles WHERE is_active = 1").Scan(&name, &model, &temp)
	if err == nil {
		fmt.Printf("- Name: %s\n- Model: %s\n- Temperature: %.2f\n\n", name, model, temp)
	}

	// 3. User Behavior
	fmt.Println("## 3. Most Active Users")
	rows, _ := db.Query("SELECT author, COUNT(*) as count FROM messages WHERE is_bot = 0 GROUP BY author ORDER BY count DESC LIMIT 10")
	for rows.Next() {
		var author string
		var count int
		rows.Scan(&author, &count)
		fmt.Printf("- %s: %d messages\n", author, count)
	}
	rows.Close()
	fmt.Println()

	// 4. Content Analysis
	fmt.Println("## 4. Topic Analysis (Keywords)")
	rows, _ = db.Query("SELECT content FROM messages WHERE is_bot = 0")
	wordCounts := make(map[string]int)
	stopWords := map[string]bool{"a": true, "o": true, "e": true, "que": true, "do": true, "da": true, "de": true, "um": true, "uma": true, "com": true, "em": true, "no": true, "na": true, "para": true, "por": true, "é": true, "não": true, "mais": true}
	for rows.Next() {
		var content string
		rows.Scan(&content)
		words := strings.Fields(strings.ToLower(content))
		for _, w := range words {
			w = strings.Trim(w, ".,!?;:()\"")
			if len(w) > 4 && !stopWords[w] && !strings.Contains(w, "<@") {
				wordCounts[w]++
			}
		}
	}
	rows.Close()

	type pair struct {
		w string
		c int
	}
	var pairs []pair
	for w, c := range wordCounts {
		pairs = append(pairs, pair{w, c})
	}
	// Sort top 20
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].c < pairs[j].c {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
		if i == 20 { break }
	}
	for i := 0; i < 20 && i < len(pairs); i++ {
		fmt.Printf("- %s (%d)\n", pairs[i].w, pairs[i].c)
	}
	fmt.Println()

	// 5. Bot Interaction Audit
	fmt.Println("## 5. Interaction Audit (Last 10 Responses)")
	rows, _ = db.Query(`
		SELECT m2.author, m2.content, m1.content 
		FROM messages m1
		JOIN messages m2 ON m2.id = (SELECT MAX(id) FROM messages WHERE id < m1.id AND is_bot = 0)
		WHERE m1.is_bot = 1
		ORDER BY m1.id DESC
		LIMIT 10
	`)
	for rows.Next() {
		var uAuthor, uContent, bContent string
		rows.Scan(&uAuthor, &uContent, &bContent)
		fmt.Printf("> **User (%s)**: %s\n", uAuthor, uContent)
		fmt.Printf("> **Bot**: %s\n\n", bContent)
	}
	rows.Close()
}
