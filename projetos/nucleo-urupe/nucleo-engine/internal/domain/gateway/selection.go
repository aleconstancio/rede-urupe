/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"sort"
	"strings"
	"unicode"

	"nucleo-engine/internal/data/sqlite"
	"github.com/aleconstancio/talos/v2/engine/intelligence"
)

func selectEpisodicCapsules(capsules []sqlite.MemoryCapsule, keywords []string, limit int) []sqlite.MemoryCapsule {
	if len(capsules) == 0 || limit <= 0 {
		return nil
	}

	type scoredCapsule struct {
		capsule sqlite.MemoryCapsule
		score   float64
	}

	var latestDailySummary *sqlite.MemoryCapsule
	var latestSummary *sqlite.MemoryCapsule
	summaryScored := make([]scoredCapsule, 0, len(capsules))
	scored := make([]scoredCapsule, 0, len(capsules))
	total := len(capsules)
	for i, capsule := range capsules {
		capsule.Normalize()
		if sqlite.IsSummaryKind(capsule.Kind) {
			if latestSummary == nil || capsule.CreatedAt.After(latestSummary.CreatedAt) {
				copy := capsule
				latestSummary = &copy
			}
			if capsule.Kind == sqlite.MemoryKindDailySummary {
				if latestDailySummary == nil || capsule.CreatedAt.After(latestDailySummary.CreatedAt) {
					copy := capsule
					latestDailySummary = &copy
				}
				continue
			}
			recencyBoost := float64(i+1) / float64(total+1)
			keywordHits := capsuleKeywordHits(capsule, keywords)
			score := float64(keywordHits*4) + recencyBoost + 0.5
			if capsule.Kind == sqlite.MemoryKindCategorySummary {
				score += 0.25
			}
			summaryScored = append(summaryScored, scoredCapsule{capsule: capsule, score: score})
			continue
		}
		if capsule.Kind == sqlite.MemoryKindCaptureSlice {
			continue
		}

		recencyBoost := float64(i+1) / float64(total+1)
		keywordHits := capsuleKeywordHits(capsule, keywords)
		openLoopBoost := 0.0
		if len(capsule.OpenLoops) > 0 {
			openLoopBoost = 1.5
		}
		score := float64(keywordHits*3) + openLoopBoost + recencyBoost
		scored = append(scored, scoredCapsule{capsule: capsule, score: score})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return scored[i].capsule.CreatedAt.After(scored[j].capsule.CreatedAt)
		}
		return scored[i].score > scored[j].score
	})

	selected := make([]sqlite.MemoryCapsule, 0, limit)
	if latestDailySummary != nil {
		selected = append(selected, *latestDailySummary)
	} else if latestSummary != nil {
		selected = append(selected, *latestSummary)
	}
	sort.SliceStable(summaryScored, func(i, j int) bool {
		if summaryScored[i].score == summaryScored[j].score {
			return summaryScored[i].capsule.CreatedAt.After(summaryScored[j].capsule.CreatedAt)
		}
		return summaryScored[i].score > summaryScored[j].score
	})
	for _, item := range summaryScored {
		if len(selected) >= limit {
			break
		}
		selected = append(selected, item.capsule)
	}
	for _, item := range scored {
		if len(selected) >= limit {
			break
		}
		selected = append(selected, item.capsule)
	}

	sort.SliceStable(selected, func(i, j int) bool {
		if selected[i].SourceEndRowID == selected[j].SourceEndRowID {
			return selected[i].CreatedAt.Before(selected[j].CreatedAt)
		}
		if selected[i].SourceEndRowID == 0 || selected[j].SourceEndRowID == 0 {
			return selected[i].CreatedAt.Before(selected[j].CreatedAt)
		}
		return selected[i].SourceEndRowID < selected[j].SourceEndRowID
	})
	return selected
}

func collectOpenLoops(capsules []sqlite.MemoryCapsule, limit int) []sqlite.OpenLoop {
	if len(capsules) == 0 || limit <= 0 {
		return nil
	}

	seen := make(map[string]struct{}, limit)
	loops := make([]sqlite.OpenLoop, 0, limit)
	for _, capsule := range capsules {
		capsule.Normalize()
		for _, loop := range capsule.OpenLoops {
			label := strings.TrimSpace(loop.Label)
			if label == "" {
				continue
			}
			key := strings.ToLower(label)
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			loops = append(loops, loop)
			if len(loops) >= limit {
				return loops
			}
		}
	}
	return loops
}

func capsuleKeywordHits(capsule sqlite.MemoryCapsule, keywords []string) int {
	if len(keywords) == 0 {
		return 0
	}

	searchBlob := capsule.SearchBlob()
	hits := 0
	for _, kw := range keywords {
		kw = strings.ToLower(strings.TrimSpace(kw))
		if kw == "" {
			continue
		}
		if strings.Contains(searchBlob, kw) {
			hits++
		}
	}
	return hits
}

func extractKeywords(messages []sqlite.Message) []string {
	if len(messages) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, 6)
	out := make([]string, 0, 6)

	for _, msg := range messages {
		for _, token := range splitKeywordTokens(strings.ToLower(msg.Content)) {
			if len(token) < 5 {
				continue
			}
			if strings.ContainsAny(token, "0123456789") {
				continue
			}
			if strings.HasPrefix(token, "http") || strings.HasPrefix(token, "www") || strings.HasPrefix(token, "@") {
				continue
			}
			// Use shared intelligence stopwords
			if intelligence.IsStopword(token) {
				continue
			}
			if _, exists := seen[token]; exists {
				continue
			}
			seen[token] = struct{}{}
			out = append(out, token)
			if len(out) >= 6 {
				return out
			}
		}
	}

	return out
}

func splitKeywordTokens(input string) []string {
	return strings.FieldsFunc(input, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}
