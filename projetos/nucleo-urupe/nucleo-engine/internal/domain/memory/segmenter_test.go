package memory

import (
	"testing"
	"time"

	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/pkg/timeutil"
)

func TestPlanCapsuleSegmentsSplitsOnSilenceGap(t *testing.T) {
	base := time.Date(2026, 5, 8, 10, 0, 0, 0, timeutil.BrasiliaLocation)
	messages := []sqlite.Message{
		testMsg(11, base.Add(0*time.Minute), "a", "alice", "question", "sql question"),
		testMsg(12, base.Add(2*time.Minute), "b", "bob", "question", "sql answer"),
		testMsg(13, base.Add(4*time.Minute), "a", "alice", "question", "follow up"),
		testMsg(14, base.Add(35*time.Minute), "c", "carol", "build", "deploy issue"),
		testMsg(15, base.Add(37*time.Minute), "d", "dan", "build", "trying rollback"),
		testMsg(16, base.Add(39*time.Minute), "c", "carol", "build", "works now"),
	}

	cfg := defaultCapsuleSegmenterConfig()
	cfg.MinMessages = 8
	cfg.MinClosedMessages = 3
	cfg.MaxSegmentsPerRun = 1

	segments := planCapsuleSegments(messages, base.Add(40*time.Minute), cfg)
	if len(segments) != 1 {
		t.Fatalf("expected 1 newest closed segment, got %d", len(segments))
	}
	if segments[0].StartExclusiveRowID != 13 || segments[0].EndInclusiveRowID != 16 {
		t.Fatalf("unexpected segment span: %+v", segments[0])
	}
	if segments[0].Reason == "" || segments[0].Reason[:11] != "silence_gap" {
		t.Fatalf("expected silence gap reason, got %q", segments[0].Reason)
	}
}

func TestPlanCapsuleSegmentsSplitsOnBudgetAndCarriesCursor(t *testing.T) {
	base := time.Date(2026, 5, 8, 12, 0, 0, 0, timeutil.BrasiliaLocation)
	messages := make([]sqlite.Message, 0, 26)
	for i := 0; i < 26; i++ {
		messages = append(messages, testMsg(int64(101+i), base.Add(time.Duration(i)*time.Minute), "a", "alice", "question", "long running discussion block"))
	}

	cfg := defaultCapsuleSegmenterConfig()
	cfg.MaxMessages = 10
	cfg.MinMessages = 8
	cfg.MinClosedMessages = 3
	cfg.MaxSegmentsPerRun = 3

	segments := planCapsuleSegments(messages, base.Add(5*time.Hour), cfg)
	if len(segments) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segments))
	}
	if segments[0].StartExclusiveRowID != 116 || segments[0].EndInclusiveRowID != 126 {
		t.Fatalf("unexpected newest segment: %+v", segments[0])
	}
	if segments[1].StartExclusiveRowID != 106 || segments[1].EndInclusiveRowID != 116 {
		t.Fatalf("unexpected second-newest segment: %+v", segments[1])
	}
	if segments[2].StartExclusiveRowID != 100 || segments[2].EndInclusiveRowID != 106 {
		t.Fatalf("unexpected oldest segment: %+v", segments[2])
	}
}

func TestPlanCapsuleSegmentsUsesCategoryShiftWithParticipantChange(t *testing.T) {
	base := time.Date(2026, 5, 8, 14, 0, 0, 0, timeutil.BrasiliaLocation)
	messages := []sqlite.Message{
		testMsg(21, base.Add(0*time.Minute), "a", "alice", "sql", "query tuning"),
		testMsg(22, base.Add(1*time.Minute), "b", "bob", "sql", "index idea"),
		testMsg(23, base.Add(2*time.Minute), "a", "alice", "sql", "planner weirdness"),
		testMsg(24, base.Add(3*time.Minute), "b", "bob", "sql", "maybe join order"),
		testMsg(25, base.Add(4*time.Minute), "a", "alice", "sql", "explain analyze"),
		testMsg(26, base.Add(5*time.Minute), "b", "bob", "sql", "looks better"),
		testMsg(27, base.Add(12*time.Minute), "c", "carol", "frontend", "new UI branch"),
		testMsg(28, base.Add(13*time.Minute), "d", "dan", "frontend", "mobile bug"),
		testMsg(29, base.Add(14*time.Minute), "c", "carol", "frontend", "fixed spacing"),
	}

	cfg := defaultCapsuleSegmenterConfig()
	cfg.MinMessages = 8
	cfg.MinClosedMessages = 3
	cfg.TopicShiftMinMessages = 6

	segments := planCapsuleSegments(messages, base.Add(3*time.Hour), cfg)
	if len(segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(segments))
	}
	if segments[0].Reason != "topic_shift:sql->frontend" {
		t.Fatalf("expected topic shift split, got %q", segments[0].Reason)
	}
	if segments[0].StartExclusiveRowID != 26 || segments[0].EndInclusiveRowID != 29 {
		t.Fatalf("unexpected newest topic-shift segment: %+v", segments[0])
	}
}

func TestPlanCapsuleSegmentsFlushesLowVolumeTailWhenStale(t *testing.T) {
	base := time.Date(2026, 5, 8, 18, 0, 0, 0, timeutil.BrasiliaLocation)
	messages := []sqlite.Message{
		testMsg(31, base.Add(0*time.Minute), "a", "alice", "question", "quick thought"),
		testMsg(32, base.Add(12*time.Minute), "b", "bob", "question", "another thought"),
		testMsg(33, base.Add(18*time.Minute), "a", "alice", "question", "we should revisit"),
	}

	cfg := defaultCapsuleSegmenterConfig()
	cfg.MinMessages = 8
	cfg.MinClosedMessages = 3

	segments := planCapsuleSegments(messages, base.Add(3*time.Hour), cfg)
	if len(segments) != 1 {
		t.Fatalf("expected stale tail to flush, got %d segments", len(segments))
	}
	if segments[0].Reason != "tail" {
		t.Fatalf("expected tail segment, got %q", segments[0].Reason)
	}
}

func TestSegmentBoundarySkipsTimestampHeuristicsWhenOrderingIsBad(t *testing.T) {
	cfg := defaultCapsuleSegmenterConfig()
	candidate := []sqlite.Message{
		testMsg(41, time.Date(2026, 5, 8, 12, 10, 0, 0, timeutil.BrasiliaLocation), "a", "alice", "sql", "latest first"),
		testMsg(42, time.Date(2026, 5, 8, 12, 5, 0, 0, timeutil.BrasiliaLocation), "b", "bob", "sql", "older second"),
		testMsg(43, time.Date(2026, 5, 8, 12, 1, 0, 0, timeutil.BrasiliaLocation), "a", "alice", "sql", "older third"),
		testMsg(44, time.Date(2026, 5, 8, 11, 58, 0, 0, timeutil.BrasiliaLocation), "b", "bob", "sql", "older fourth"),
		testMsg(45, time.Date(2026, 5, 8, 11, 55, 0, 0, timeutil.BrasiliaLocation), "a", "alice", "sql", "older fifth"),
		testMsg(46, time.Date(2026, 5, 8, 11, 50, 0, 0, timeutil.BrasiliaLocation), "b", "bob", "sql", "older sixth"),
	}
	current := candidate[len(candidate)-1]
	next := testMsg(47, time.Date(2026, 5, 8, 11, 30, 0, 0, timeutil.BrasiliaLocation), "c", "carol", "frontend", "looks like a gap")

	if !significantCategoryShift(current.Category, next.Category) {
		t.Fatal("expected category shift precondition")
	}
	if _, ok := segmentBoundaryReason(candidate, current, next, cfg); ok {
		t.Fatal("expected boundary heuristic to be disabled for badly ordered timestamps")
	}
}

func TestPlanCapsuleSegmentsDoesNotSkipFreshNewestTail(t *testing.T) {
	base := time.Date(2026, 5, 8, 20, 0, 0, 0, timeutil.BrasiliaLocation)
	messages := []sqlite.Message{
		testMsg(51, base.Add(0*time.Minute), "a", "alice", "sql", "older one"),
		testMsg(52, base.Add(1*time.Minute), "b", "bob", "sql", "older two"),
		testMsg(53, base.Add(2*time.Minute), "a", "alice", "sql", "older three"),
		testMsg(54, base.Add(30*time.Minute), "c", "carol", "frontend", "fresh newest"),
	}

	cfg := defaultCapsuleSegmenterConfig()
	cfg.MinMessages = 8
	cfg.MinClosedMessages = 3
	cfg.SingleMessageStaleAfter = 12 * time.Hour
	cfg.StaleAfter = 45 * time.Minute

	segments := planCapsuleSegments(messages, base.Add(31*time.Minute), cfg)
	if len(segments) != 0 {
		t.Fatalf("expected no segment while freshest tail is not ready, got %+v", segments)
	}
}

func testMsg(id int64, ts time.Time, authorID, author, category, content string) sqlite.Message {
	return sqlite.Message{
		ID:        id,
		Timestamp: ts,
		AuthorID:  authorID,
		Author:    author,
		Category:  category,
		Content:   content,
		ChannelID: "chan",
	}
}
