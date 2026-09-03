/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Model struct {
	ID string `json:"id"`
}

func main() {
	apiKey := strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("LLM_API_KEY"))
	}
	if apiKey == "" {
		fmt.Println("OPENROUTER_API_KEY or LLM_API_KEY is required")
		return
	}

	req, err := http.NewRequest("GET", "https://openrouter.ai/api/v1/models", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var res struct {
		Data []Model `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		return
	}

	fmt.Println("Available free models on OpenRouter:")
	for _, m := range res.Data {
		if strings.HasSuffix(m.ID, ":free") || m.ID == "openrouter/free" {
			fmt.Println(m.ID)
		}
	}
}
