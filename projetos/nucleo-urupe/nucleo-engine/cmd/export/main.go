/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"nucleo-engine/internal/config"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type TrainingExample struct {
	Messages []TrainingMessage `json:"messages"`
}

type TrainingMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
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

	fmt.Println("Exporting data for fine-tuning...")

	// Get all messages ordered by time
	rows, err := db.Query(`
		SELECT author, content, timestamp, is_bot, channel_id 
		FROM messages 
		ORDER BY channel_id, timestamp ASC
	`)
	if err != nil {
		log.Fatalf("Failed to query messages: %v", err)
	}
	defer rows.Close()

	file, err := os.Create("training_data.jsonl")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	var currentConversation []TrainingMessage
	var lastTimestamp time.Time
	var lastChannel string
	
	const sessionGap = 15 * time.Minute
	count := 0

	for rows.Next() {
		var author, content, tsStr string
		var isBot bool
		var channelID string
		rows.Scan(&author, &content, &tsStr, &isBot, &channelID)

		ts, _ := time.Parse("2006-01-02 15:04:05", tsStr)

		// Start new conversation if:
		// 1. Channel changed
		// 2. Large time gap
		if (lastChannel != "" && channelID != lastChannel) || (!lastTimestamp.IsZero() && ts.Sub(lastTimestamp) > sessionGap) {
			if len(currentConversation) > 1 {
				saveConversation(file, currentConversation)
				count++
			}
			currentConversation = nil
		}

		role := "user"
		if isBot {
			role = "assistant"
		}

		currentConversation = append(currentConversation, TrainingMessage{
			Role:    role,
			Content: fmt.Sprintf("[%s]: %s", author, content),
			Name:    sanitizeName(author),
		})

		lastTimestamp = ts
		lastChannel = channelID

		// Limit conversation size to avoid context overflow in training
		if len(currentConversation) > 20 {
			saveConversation(file, currentConversation)
			count++
			currentConversation = nil
		}
	}

	// Final conversation
	if len(currentConversation) > 1 {
		saveConversation(file, currentConversation)
		count++
	}

	fmt.Printf("Successfully exported %d conversation sessions to training_data.jsonl\n", count)
}

func saveConversation(file *os.File, msgs []TrainingMessage) {
	example := TrainingExample{Messages: msgs}
	data, _ := json.Marshal(example)
	file.Write(data)
	file.WriteString("\n")
}

func sanitizeName(name string) string {
	// OpenAI/Llama names must match ^[a-zA-Z0-9_-]{1,64}$
	out := ""
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			out += string(r)
		}
	}
	if len(out) > 64 {
		out = out[:64]
	}
	return out
}
