/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken      string
	OpenRouterAPIKey  string
	OpenRouterModel   string
	ResponseFallbacks []string
	SummaryFallbacks  []string
	OpenRouterSiteURL string
	OpenRouterTitle   string

	// Gemini Native Config
	GeminiAPIKeyFree   string
	GeminiAPIKey       string
	GeminiAPIKey2      string
	GroqAPIKey         string
	MetacognitionModel string // legacy fallback for memory work
	SynthesisModel     string // legacy fallback for reply generation
	VisionModel        string // e.g. gemini-3-flash
	GatewayGateModel   string // e.g. gemini-3.1-flash-lite
	GatewayReplyModel  string // e.g. gemini-3.1-flash
	MemoryModel        string // e.g. gemini-3.1-flash-lite

	TargetGuildID       string
	TargetChannelID     string
	TargetChannelName   string
	MultiChannel        bool
	DatabasePath        string
	DashboardEnabled    bool
	DashboardPort       string
	ResponseTemperature float64
	MaxOutputTokens     int
	AdminAPIKey         string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DiscordToken:      strings.TrimSpace(os.Getenv("DISCORD_TOKEN")),
		OpenRouterAPIKey:  firstNonEmpty("OPENROUTER_API_KEY", "LLM_API_KEY"),
		OpenRouterModel:   firstNonEmpty("OPENROUTER_MODEL", "LLM_MODEL"),
		ResponseFallbacks: defaultCSVKeys([]string{"OPENROUTER_RESPONSE_FALLBACK_MODELS", "OPENROUTER_FALLBACK_MODELS"}, defaultResponseFallbackModels()),
		SummaryFallbacks:  defaultCSVKeys([]string{"OPENROUTER_SUMMARY_FALLBACK_MODELS"}, defaultSummaryFallbackModels()),
		OpenRouterSiteURL: defaultString("OPENROUTER_SITE_URL", "http://localhost"),
		OpenRouterTitle:   defaultString("OPENROUTER_APP_NAME", "Núcleo Urupê"),

		GeminiAPIKeyFree:   strings.TrimSpace(os.Getenv("GEMINI_API_KEY_FREE")),
		GeminiAPIKey:       strings.TrimSpace(os.Getenv("GEMINI_API_KEY")),
		GeminiAPIKey2:      strings.TrimSpace(os.Getenv("GEMINI_API_KEY_2")),
		GroqAPIKey:         strings.TrimSpace(os.Getenv("GROQ_API_KEY")),
		MetacognitionModel: defaultString("GEMINI_METACOGNITION_MODEL", "gemini-3.1-flash-lite"),
		SynthesisModel:     defaultString("GEMINI_SYNTHESIS_MODEL", "gemini-3.1-flash-lite"),
		VisionModel:        defaultString("GEMINI_VISION_MODEL", "gemini-3.0-flash-preview"),
		GatewayGateModel:   defaultString("GEMINI_GATEWAY_GATE_MODEL", ""),
		GatewayReplyModel:  defaultString("GEMINI_GATEWAY_REPLY_MODEL", ""),
		MemoryModel:        defaultString("GEMINI_MEMORY_MODEL", ""),

		TargetGuildID:       strings.TrimSpace(os.Getenv("DISCORD_GUILD_ID")),
		TargetChannelID:     strings.TrimSpace(os.Getenv("DISCORD_CHANNEL_ID")),
		DatabasePath:        defaultString("SQLITE_PATH", "data/nucleo.db"),
		DashboardEnabled:    defaultBool("DASHBOARD_ENABLED", true),
		DashboardPort:       defaultString("DASHBOARD_PORT", "9393"),
		ResponseTemperature: defaultFloat("LLM_TEMPERATURE", 0.4),
		MaxOutputTokens:     defaultInt("LLM_MAX_OUTPUT_TOKENS", 3200),
		AdminAPIKey:         strings.TrimSpace(os.Getenv("ADMIN_API_KEY")),
	}

	if cfg.DiscordToken == "" {
		return Config{}, fmt.Errorf("DISCORD_TOKEN environment variable is required")
	}
	if cfg.OpenRouterAPIKey == "" && cfg.GeminiAPIKey == "" {
		return Config{}, fmt.Errorf("either OPENROUTER_API_KEY or GEMINI_API_KEY environment variable is required")
	}
	cfg.MultiChannel = cfg.TargetChannelID == ""
	if cfg.OpenRouterModel == "" {
		cfg.OpenRouterModel = "meta-llama/llama-3.3-70b-instruct:free"
	}
	if cfg.GatewayGateModel == "" {
		cfg.GatewayGateModel = firstNonEmpty("GEMINI_METACOGNITION_MODEL")
	}
	if cfg.GatewayGateModel == "" {
		cfg.GatewayGateModel = "gemini-3.1-flash-lite"
	}
	if cfg.GatewayReplyModel == "" {
		cfg.GatewayReplyModel = firstNonEmpty("GEMINI_SYNTHESIS_MODEL")
	}
	if cfg.GatewayReplyModel == "" {
		cfg.GatewayReplyModel = "gemini-3.1-flash-lite"
	}
	if cfg.MemoryModel == "" {
		cfg.MemoryModel = firstNonEmpty("GEMINI_METACOGNITION_MODEL")
	}
	if cfg.MemoryModel == "" {
		cfg.MemoryModel = "gemini-3.1-flash-lite"
	}
	if cfg.MaxOutputTokens < 64 {
		cfg.MaxOutputTokens = 64
	}

	return cfg, nil
}

func firstNonEmpty(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func defaultString(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func defaultBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func defaultInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func defaultFloat(key string, fallback float64) float64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func defaultCSVKeys(keys []string, fallback []string) []string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value == "" {
			continue
		}

		parts := strings.Split(value, ",")
		out := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
		if len(out) > 0 {
			return out
		}
	}
	return append([]string(nil), fallback...)
}

func defaultResponseFallbackModels() []string {
	return []string{
		"tencent/hy3-preview:free",
		"nvidia/nemotron-3-super:free",
		"inclusionai/ling-2.6-1t:free",
		"nous/hermes-3-405b-instruct:free",
		"openai/gpt-oss-120b:free",
		"qwen/qwen3-coder-480b-a35b:free",
		"poolside/laguna-m.1:free",
		"poolside/laguna-xs.2:free",
		"owl/owl-alpha:free",
		"google/gemma-4-31b:free",
		"google/gemma-4-26b-a4b:free",
		"meta/llama-3.3-70b-instruct:free",
		"liquidai/lfm2.5-1.2b-thinking:free",
		"liquidai/lfm2.5-1.2b-instruct:free",
		"venice/uncensored:free",
		"minimax/minimax-m2.5:free",
		"nvidia/nemotron-3-nano-30b-a3b:free",
		"qwen/qwen3-next-80b-a3b-instruct:free",
		"nvidia/nemotron-nano-9b-v2:free",
		"openai/gpt-oss-20b:free",
		"z-ai/glm-4.5-air:free",
		"google/gemma-3n-2b:free",
		"google/gemma-3n-4b:free",
		"google/gemma-3-4b:free",
		"google/gemma-3-12b:free",
		"google/gemma-3-27b:free",
		"meta/llama-3.2-3b-instruct:free",
	}
}

func defaultSummaryFallbackModels() []string {
	return []string{
		"nous/hermes-3-405b-instruct:free",
		"openai/gpt-oss-120b:free",
		"nvidia/nemotron-3-super:free",
		"inclusionai/ling-2.6-1t:free",
		"tencent/hy3-preview:free",
		"qwen/qwen3-coder-480b-a35b:free",
		"poolside/laguna-m.1:free",
		"google/gemma-4-31b:free",
		"google/gemma-4-26b-a4b:free",
		"meta/llama-3.3-70b-instruct:free",
		"minimax/minimax-m2.5:free",
		"qwen/qwen3-next-80b-a3b-instruct:free",
		"openai/gpt-oss-20b:free",
		"z-ai/glm-4.5-air:free",
		"google/gemma-3-27b:free",
		"google/gemma-3-12b:free",
		"google/gemma-3-4b:free",
		"meta/llama-3.2-3b-instruct:free",
	}
}
