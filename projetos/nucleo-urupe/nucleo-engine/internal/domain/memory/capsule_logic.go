/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package memory

import (
	"fmt"
	"strings"

	"nucleo-engine/internal/data/sqlite"
	"github.com/aleconstancio/talos/v2/engine/intelligence"
	"nucleo-engine/internal/pkg/timeutil"
)

func hydrateCapsuleFromSegment(capsule *sqlite.MemoryCapsule, segment capsuleSegment) {
	startMsg := segment.Messages[0]
	endMsg := segment.Messages[len(segment.Messages)-1]

	capsule.DayDate = timeutil.InBrasilia(endMsg.Timestamp).Format("2006-01-02")
	capsule.TimeSpan = fmt.Sprintf("%s-%s", timeutil.InBrasilia(startMsg.Timestamp).Format("15:04"), timeutil.InBrasilia(endMsg.Timestamp).Format("15:04"))
	capsule.SourceStartRowID = startMsg.ID
	capsule.SourceEndRowID = endMsg.ID
	capsule.SourceMessageCount = len(segment.Messages)
	capsule.IsMerged = false
	capsule.Category = dominantMessageCategory(segment.Messages)

	sanitizeCapsuleParticipants(capsule, segment.Messages)
	sanitizeCapsuleOpenLoops(capsule, segment.Messages)
	sparsifyCapsule(capsule)
	capsule.Kind = determineCapsuleKind(*capsule, segment)
	if strings.TrimSpace(capsule.MainTopic) == "" {
		capsule.MainTopic = fallbackMainTopic(segment.Messages)
	}
	if strings.TrimSpace(capsule.Mood) == "" {
		capsule.Mood = "neutral"
	}
	capsule.Normalize()
}

func sanitizeCapsuleParticipants(capsule *sqlite.MemoryCapsule, messages []sqlite.Message) {
	allowed := make(map[string]struct{}, len(messages))
	fallback := make([]string, 0, len(messages))
	for _, msg := range messages {
		name := strings.TrimSpace(msg.Author)
		if name == "" {
			continue
		}
		if _, exists := allowed[name]; exists {
			continue
		}
		allowed[name] = struct{}{}
		fallback = append(fallback, name)
	}

	filtered := make([]string, 0, len(capsule.Participants))
	seen := make(map[string]struct{}, len(capsule.Participants))
	for _, participant := range capsule.Participants {
		participant = strings.TrimSpace(participant)
		if participant == "" {
			continue
		}
		if _, ok := allowed[participant]; !ok {
			continue
		}
		if _, exists := seen[participant]; exists {
			continue
		}
		seen[participant] = struct{}{}
		filtered = append(filtered, participant)
	}

	if len(filtered) == 0 {
		filtered = fallback
	}
	if len(filtered) > 8 {
		filtered = filtered[:8]
	}
	capsule.Participants = filtered
}

func sanitizeCapsuleOpenLoops(capsule *sqlite.MemoryCapsule, messages []sqlite.Message) {
	allowedOwners := make(map[string]struct{}, len(messages))
	for _, msg := range messages {
		if author := strings.TrimSpace(msg.Author); author != "" {
			allowedOwners[author] = struct{}{}
		}
	}

	loops := make([]sqlite.OpenLoop, 0, len(capsule.OpenLoops))
	for _, loop := range capsule.OpenLoops {
		loop.Label = strings.TrimSpace(loop.Label)
		loop.Owner = strings.TrimSpace(loop.Owner)
		loop.NextStep = strings.TrimSpace(loop.NextStep)
		loop.Status = strings.TrimSpace(loop.Status)
		if loop.Label == "" {
			continue
		}
		if loop.Owner != "" {
			if _, ok := allowedOwners[loop.Owner]; !ok {
				loop.Owner = ""
			}
		}
		loops = append(loops, loop)
	}
	capsule.OpenLoops = loops
}

func sparsifyCapsule(capsule *sqlite.MemoryCapsule) {
	capsule.TypedFacts.Events = trimStringList(capsule.TypedFacts.Events, 3)
	capsule.TypedFacts.Traits = trimStringList(capsule.TypedFacts.Traits, 2)
	capsule.TypedFacts.Tensions = trimStringList(capsule.TypedFacts.Tensions, 2)
	capsule.TypedFacts.Callbacks = trimStringList(capsule.TypedFacts.Callbacks, 1)
	capsule.OpenLoops = trimOpenLoops(capsule.OpenLoops, 3)
}

func trimStringList(items []string, limit int) []string {
	if limit <= 0 {
		return nil
	}

	out := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		key := strings.ToLower(item)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, item)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func trimOpenLoops(loops []sqlite.OpenLoop, limit int) []sqlite.OpenLoop {
	if limit <= 0 {
		return nil
	}

	out := make([]sqlite.OpenLoop, 0, len(loops))
	seen := make(map[string]struct{}, len(loops))
	for _, loop := range loops {
		label := strings.TrimSpace(loop.Label)
		if label == "" {
			continue
		}
		key := strings.ToLower(label)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		loop.Label = label
		out = append(out, loop)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func determineCapsuleKind(capsule sqlite.MemoryCapsule, segment capsuleSegment) string {
	score := capsuleSignalScore(capsule)
	durableSignals := len(capsule.TypedFacts.Traits) + len(capsule.TypedFacts.Tensions) + len(capsule.OpenLoops)
	eventSignals := len(capsule.TypedFacts.Events) + len(capsule.TypedFacts.Callbacks)

	// Check for high-signal categories in the segment
	hasHighSignalCategory := false
	for _, msg := range segment.Messages {
		cat := msg.Category
		if cat == string(intelligence.QuickInquiry) || cat == string(intelligence.SocialInquiry) || cat == string(intelligence.ConceptualInquiry) {
			hasHighSignalCategory = true
			break
		}
	}

	switch {
	case durableSignals >= 2:
		return sqlite.MemoryKindEpisodeSnapshot
	case durableSignals >= 1 && segment.MessageCount >= 24:
		return sqlite.MemoryKindEpisodeSnapshot
	case hasHighSignalCategory && segment.MessageCount >= 4:
		return sqlite.MemoryKindEpisodeSnapshot
	case score >= 6:
		return sqlite.MemoryKindEpisodeSnapshot
	case segment.MessageCount >= 90 && eventSignals >= 2:
		return sqlite.MemoryKindEpisodeSnapshot
	default:
		return sqlite.MemoryKindCaptureSlice
	}
}

func capsuleSignalScore(capsule sqlite.MemoryCapsule) int {
	return len(capsule.TypedFacts.Events) +
		len(capsule.TypedFacts.Callbacks) +
		(2 * len(capsule.OpenLoops)) +
		(3 * len(capsule.TypedFacts.Traits)) +
		(3 * len(capsule.TypedFacts.Tensions))
}

func fallbackMainTopic(messages []sqlite.Message) string {
	category := dominantMessageCategory(messages)
	if category == "" || category == "general" || category == "uncategorized" {
		return "mixed conversation segment"
	}
	return category
}

func dominantMessageCategory(messages []sqlite.Message) string {
	counts := make(map[string]int)
	order := make([]string, 0, len(messages))
	for _, msg := range messages {
		category := msg.Category
		if category == "" {
			category = "uncategorized"
		}
		if _, seen := counts[category]; !seen {
			order = append(order, category)
		}
		counts[category]++
	}

	best := "general"
	bestCount := -1
	for _, category := range order {
		if counts[category] > bestCount {
			best = category
			bestCount = counts[category]
		}
	}
	return best
}
