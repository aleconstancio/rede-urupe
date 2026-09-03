/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */

// Package moderation provides rule-based spam, toxicity, and personal attack detection.
package moderation

import (
	"strings"
	"time"

	"github.com/aleconstancio/talos/v2/engine/intelligence"
)

type RuleHit struct {
	Rule    string
	Reason  string
	Score   float64
	Action  string // "warn", "delete", "timeout"
	ActionDuration time.Duration
}

type ModContext struct {
	Content   string
	AuthorID  string
	ChannelID string
	GuildID   string
	MessageID string
	Recent    []string
}

func CheckMessage(ctx ModContext, tone intelligence.ToneProfile) []RuleHit {
	var hits []RuleHit

	if hit := checkSpam(ctx); hit != nil {
		hits = append(hits, *hit)
	}

	if hit := checkToxicity(tone); hit != nil {
		hits = append(hits, *hit)
	}

	if hit := checkPersonalAttack(ctx); hit != nil {
		hits = append(hits, *hit)
	}

	return hits
}

func checkSpam(ctx ModContext) *RuleHit {
	if len(ctx.Recent) < 5 {
		return nil
	}

	recent := ctx.Recent
	if len(recent) > 10 {
		recent = recent[len(recent)-10:]
	}

	exactDups := 0
	last := strings.TrimSpace(ctx.Content)
	if last == "" {
		return nil
	}
	for _, m := range recent {
		if strings.TrimSpace(m) == last {
			exactDups++
		}
	}

	if exactDups >= 3 {
		return &RuleHit{
			Rule:    "spam_duplicate",
			Reason:  "conteudo repetido detectado",
			Score:   0.8,
			Action:  "warn",
			ActionDuration: 0,
		}
	}

	shortCount := 0
	for _, m := range recent {
		if len(strings.Fields(m)) <= 3 {
			shortCount++
		}
	}
	if shortCount >= 8 {
		return &RuleHit{
			Rule:    "spam_flood",
			Reason:  "mensagens curtas em massa",
			Score:   0.7,
			Action:  "warn",
		}
	}

	return nil
}

func checkToxicity(tone intelligence.ToneProfile) *RuleHit {
	if tone.Conflict >= 0.7 {
		return &RuleHit{
			Rule:   "toxicity_conflict",
			Reason: "tom de conflito elevado",
			Score:  tone.Conflict,
			Action: "warn",
		}
	}
	if tone.Intensity >= 0.8 && tone.Conflict >= 0.5 {
		return &RuleHit{
			Rule:   "toxicity_aggressive",
			Reason: "intensidade agressiva com indicio de conflito",
			Score:  tone.Intensity * 0.7,
			Action: "warn",
		}
	}
	return nil
}

func checkPersonalAttack(ctx ModContext) *RuleHit {
	lower := strings.ToLower(ctx.Content)
	attackWords := []string{
		"burro", "idiota", "retardado", "imbecil", "otario", "otário",
		"arrombado", "filho da puta", "vai tomar no cu", "vai se foder",
		"seu lixo", "sua anta", "corno", "trouxa",
	}
	for _, w := range attackWords {
		if strings.Contains(lower, w) {
			return &RuleHit{
				Rule:   "personal_attack",
				Reason: "ataque pessoal detectado",
				Score:  0.9,
				Action: "warn",
			}
		}
	}
	return nil
}
