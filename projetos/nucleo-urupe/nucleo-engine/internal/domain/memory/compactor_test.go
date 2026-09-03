package memory

import (
	"testing"
	"time"

	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/pkg/timeutil"
)

func TestChooseSummaryKindHeuristicPrefersCategorySummaryForDominantCategory(t *testing.T) {
	capsules := []sqlite.MemoryCapsule{
		{Kind: sqlite.MemoryKindEpisodeSnapshot, Category: "memory", MainTopic: "routing", CreatedAt: time.Date(2026, 5, 8, 9, 0, 0, 0, timeutil.BrasiliaLocation)},
		{Kind: sqlite.MemoryKindEpisodeSnapshot, Category: "memory", MainTopic: "segmenter", CreatedAt: time.Date(2026, 5, 8, 10, 0, 0, 0, timeutil.BrasiliaLocation)},
		{Kind: sqlite.MemoryKindEpisodeSnapshot, Category: "memory", MainTopic: "prompt assembly", CreatedAt: time.Date(2026, 5, 8, 11, 0, 0, 0, timeutil.BrasiliaLocation)},
		{Kind: sqlite.MemoryKindEpisodeSnapshot, Category: "humor", MainTopic: "banter", CreatedAt: time.Date(2026, 5, 8, 12, 0, 0, 0, timeutil.BrasiliaLocation)},
	}

	if got := chooseSummaryKindHeuristic(capsules); got != sqlite.MemoryKindCategorySummary {
		t.Fatalf("expected category summary heuristic, got %s", got)
	}
}

func TestSelectCompactionInputKeepsLatestSummaryAcrossSummaryKinds(t *testing.T) {
	active := []sqlite.MemoryCapsule{
		{
			ID:        1,
			Kind:      sqlite.MemoryKindDailySummary,
			DayDate:   "2026-05-08",
			TimeSpan:  "09:00-09:30",
			EpisodeID: "2026-05-08:daily_summary:1-20",
			CreatedAt: time.Date(2026, 5, 8, 9, 30, 0, 0, timeutil.BrasiliaLocation),
		},
		{
			ID:        2,
			Kind:      sqlite.MemoryKindTopicSummary,
			DayDate:   "2026-05-08",
			TimeSpan:  "10:00-10:20",
			EpisodeID: "2026-05-08:topic_summary:21-32",
			CreatedAt: time.Date(2026, 5, 8, 10, 20, 0, 0, timeutil.BrasiliaLocation),
		},
		{
			ID:        3,
			Kind:      sqlite.MemoryKindEpisodeSnapshot,
			DayDate:   "2026-05-08",
			TimeSpan:  "11:00-11:20",
			EpisodeID: "2026-05-08:episode_snapshot:33-44",
			CreatedAt: time.Date(2026, 5, 8, 11, 20, 0, 0, timeutil.BrasiliaLocation),
		},
	}

	got := selectCompactionInput(active)
	if len(got) != 2 {
		t.Fatalf("expected one latest summary plus one snapshot, got %d: %+v", len(got), got)
	}
	if got[0].Kind != sqlite.MemoryKindTopicSummary {
		t.Fatalf("expected latest summary kind to win, got %+v", got)
	}
	if got[1].Kind != sqlite.MemoryKindEpisodeSnapshot {
		t.Fatalf("expected snapshot to remain eligible, got %+v", got)
	}
}
