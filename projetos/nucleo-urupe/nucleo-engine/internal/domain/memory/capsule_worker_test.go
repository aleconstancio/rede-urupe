package memory

import (
	"testing"
	"time"

	"nucleo-engine/internal/data/sqlite"
	"github.com/aleconstancio/talos/v2/engine/intelligence"
	"nucleo-engine/internal/pkg/timeutil"
)

func TestDetermineCapsuleKindPrefersDurableSignals(t *testing.T) {
	segment := capsuleSegment{
		MessageCount: 36,
		Messages: []sqlite.Message{
			{ID: 10, Timestamp: time.Date(2026, 5, 8, 12, 0, 0, 0, timeutil.BrasiliaLocation)},
			{ID: 45, Timestamp: time.Date(2026, 5, 8, 13, 0, 0, 0, timeutil.BrasiliaLocation)},
		},
	}

	capsule := sqlite.MemoryCapsule{
		TypedFacts: sqlite.TypedFacts{
			Traits: []string{"Ale prefere respostas mais concisas e com mais criterio de relevancia"},
		},
	}

	if got := determineCapsuleKind(capsule, segment); got != sqlite.MemoryKindEpisodeSnapshot {
		t.Fatalf("expected durable trait signal to produce episode snapshot, got %s", got)
	}
}

func TestDetermineCapsuleKindDemotesThinEventOnlySegments(t *testing.T) {
	segment := capsuleSegment{
		MessageCount: 34,
		Messages: []sqlite.Message{
			{ID: 100, Timestamp: time.Date(2026, 5, 8, 12, 0, 0, 0, timeutil.BrasiliaLocation)},
			{ID: 133, Timestamp: time.Date(2026, 5, 8, 13, 0, 0, 0, timeutil.BrasiliaLocation)},
		},
	}

	capsule := sqlite.MemoryCapsule{
		TypedFacts: sqlite.TypedFacts{
			Events: []string{"O grupo falou sobre um tema e seguiu em frente"},
		},
	}

	if got := determineCapsuleKind(capsule, segment); got != sqlite.MemoryKindCaptureSlice {
		t.Fatalf("expected thin event-only segment to become capture slice, got %s", got)
	}
}

func TestDetermineCapsuleKindPromotesDirectInquirySegments(t *testing.T) {
	segment := capsuleSegment{
		MessageCount: 8,
		Messages: []sqlite.Message{
			{
				ID:        200,
				Timestamp: time.Date(2026, 5, 8, 12, 0, 0, 0, timeutil.BrasiliaLocation),
				Author:    "Ale",
				Category:  string(intelligence.QuickInquiry),
				Content:   "Talos, you can give feedback on this approach?",
			},
			{
				ID:        201,
				Timestamp: time.Date(2026, 5, 8, 12, 1, 0, 0, timeutil.BrasiliaLocation),
				Author:    "Talos",
				IsBot:     true,
				Category:  string(intelligence.SocialInquiry),
				Content:   "Claro, posso olhar isso.",
			},
		},
	}

	capsule := sqlite.MemoryCapsule{
		TypedFacts: sqlite.TypedFacts{
			Events: []string{"feedback was requested"},
		},
	}

	if got := determineCapsuleKind(capsule, segment); got != sqlite.MemoryKindEpisodeSnapshot {
		t.Fatalf("expected direct inquiry segment to promote to episode snapshot, got %s", got)
	}
}

func TestSparsifyCapsuleKeepsMemoryConcise(t *testing.T) {
	capsule := sqlite.MemoryCapsule{
		TypedFacts: sqlite.TypedFacts{
			Events:    []string{"a", "b", "c", "d"},
			Traits:    []string{"t1", "t2", "t3"},
			Tensions:  []string{"x", "y", "z"},
			Callbacks: []string{"cb1", "cb2"},
		},
		OpenLoops: []sqlite.OpenLoop{
			{Label: "l1"},
			{Label: "l2"},
			{Label: "l3"},
			{Label: "l4"},
		},
	}

	sparsifyCapsule(&capsule)

	if len(capsule.TypedFacts.Events) != 3 || len(capsule.TypedFacts.Traits) != 2 || len(capsule.TypedFacts.Tensions) != 2 || len(capsule.TypedFacts.Callbacks) != 1 || len(capsule.OpenLoops) != 3 {
		t.Fatalf("unexpected sparse capsule shape: %+v", capsule)
	}
}
