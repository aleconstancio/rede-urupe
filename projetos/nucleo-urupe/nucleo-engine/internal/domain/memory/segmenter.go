package memory

import (
	"fmt"
	"strings"
	"time"

	"nucleo-engine/internal/data/sqlite"
)

type capsuleSegment struct {
	StartExclusiveRowID int64
	EndInclusiveRowID   int64
	Messages            []sqlite.Message
	Reason              string
	MessageCount        int
	ApproxChars         int
	OldestTimestamp     time.Time
	NewestTimestamp     time.Time
}

type capsuleSegmenterConfig struct {
	MinMessages             int
	MinClosedMessages       int
	MaxMessages             int
	SoftCharBudget          int
	HardCharBudget          int
	StaleAfter              time.Duration
	SilenceGap              time.Duration
	TopicShiftGap           time.Duration
	TopicShiftMinMessages   int
	TailAuthorWindow        int
	MaxSegmentsPerRun       int
	TailProbeMessages       int
	OrderingProbeWindow     int
	MaxNegativeGapRatio     float64
	SingleMessageStaleAfter time.Duration
}

func defaultCapsuleSegmenterConfig() capsuleSegmenterConfig {
	return capsuleSegmenterConfig{
		MinMessages:             30,
		MinClosedMessages:       12,
		MaxMessages:             120,
		SoftCharBudget:          16000,
		HardCharBudget:          22000,
		StaleAfter:              2 * time.Hour,
		SilenceGap:              20 * time.Minute,
		TopicShiftGap:           10 * time.Minute,
		TopicShiftMinMessages:   24,
		TailAuthorWindow:        8,
		MaxSegmentsPerRun:       2,
		TailProbeMessages:       384,
		OrderingProbeWindow:     24,
		MaxNegativeGapRatio:     0.20,
		SingleMessageStaleAfter: 24 * time.Hour,
	}
}

func planCapsuleSegments(messages []sqlite.Message, now time.Time, cfg capsuleSegmenterConfig) []capsuleSegment {
	if len(messages) == 0 {
		return nil
	}
	if cfg.MaxSegmentsPerRun <= 0 {
		cfg.MaxSegmentsPerRun = 1
	}

	segments := make([]capsuleSegment, 0, cfg.MaxSegmentsPerRun)
	end := len(messages)
	rightClosed := false

	for end > 0 && len(segments) < cfg.MaxSegmentsPerRun {
		segment, start, ok := planNewestCapsuleSegment(messages[:end], now, cfg, rightClosed)
		if !ok {
			break
		}
		segments = append(segments, segment)
		end = start
		rightClosed = true
	}
	return segments
}

func planNewestCapsuleSegment(messages []sqlite.Message, now time.Time, cfg capsuleSegmenterConfig, rightClosed bool) (capsuleSegment, int, bool) {
	if len(messages) == 0 {
		return capsuleSegment{}, 0, false
	}

	end := len(messages)
	start := end - 1

	for start > 0 {
		candidate := buildCapsuleSegment(messages[start:end], segmentStartExclusiveRowID(messages, start), "")
		if exceedsSegmentBudget(buildCapsuleSegment(messages[start-1:end], segmentStartExclusiveRowID(messages, start-1), ""), cfg) {
			candidate.Reason = "forced_budget_cut"
			if !shouldFlushSegment(candidate, now, cfg, true) {
				return capsuleSegment{}, 0, false
			}
			return candidate, start, true
		}

		if reason, ok := reverseSegmentBoundaryReason(messages[:end], start, cfg); ok {
			candidate.Reason = reason
			if shouldFlushSegment(candidate, now, cfg, true) {
				return candidate, start, true
			}
		}

		start--
	}

	tail := buildCapsuleSegment(messages[:end], segmentStartExclusiveRowID(messages, 0), "tail")
	if !shouldFlushSegment(tail, now, cfg, rightClosed) {
		return capsuleSegment{}, 0, false
	}
	return tail, 0, true
}

func buildCapsuleSegment(messages []sqlite.Message, startExclusiveRowID int64, reason string) capsuleSegment {
	if len(messages) == 0 {
		return capsuleSegment{}
	}

	approxChars := 0
	for _, msg := range messages {
		approxChars += approximateMessageChars(msg)
	}

	return capsuleSegment{
		StartExclusiveRowID: startExclusiveRowID,
		EndInclusiveRowID:   messages[len(messages)-1].ID,
		Messages:            messages,
		Reason:              reason,
		MessageCount:        len(messages),
		ApproxChars:         approxChars,
		OldestTimestamp:     messages[0].Timestamp,
		NewestTimestamp:     messages[len(messages)-1].Timestamp,
	}
}

func approximateMessageChars(msg sqlite.Message) int {
	size := len(strings.TrimSpace(msg.Content)) + 24
	size += len(msg.Attachments) * 32
	size += len(msg.Reactions) * 8
	if msg.ReplyToID != "" {
		size += 12
	}
	return size
}

func exceedsSegmentBudget(seg capsuleSegment, cfg capsuleSegmenterConfig) bool {
	return seg.MessageCount > cfg.MaxMessages || seg.ApproxChars > cfg.SoftCharBudget
}

func shouldFlushSegment(seg capsuleSegment, now time.Time, cfg capsuleSegmenterConfig, closed bool) bool {
	if seg.MessageCount == 0 {
		return false
	}
	if seg.MessageCount == 1 {
		return now.Sub(seg.OldestTimestamp) >= cfg.SingleMessageStaleAfter || seg.ApproxChars >= cfg.HardCharBudget
	}
	if seg.MessageCount >= cfg.MinMessages || seg.ApproxChars >= cfg.HardCharBudget {
		return true
	}
	if closed && seg.MessageCount >= cfg.MinClosedMessages {
		return true
	}
	return seg.MessageCount >= cfg.MinClosedMessages && now.Sub(seg.OldestTimestamp) >= cfg.StaleAfter
}

func segmentStartExclusiveRowID(messages []sqlite.Message, start int) int64 {
	if len(messages) == 0 {
		return 0
	}
	if start <= 0 {
		if messages[0].ID <= 0 {
			return 0
		}
		return messages[0].ID - 1
	}
	return messages[start-1].ID
}

func reverseSegmentBoundaryReason(messages []sqlite.Message, start int, cfg capsuleSegmenterConfig) (string, bool) {
	if start <= 0 || start >= len(messages) {
		return "", false
	}

	olderIdx := start - 1
	probeStart := boundaryProbeStart(olderIdx, cfg)
	return segmentBoundaryReason(messages[probeStart:olderIdx+1], messages[olderIdx], messages[start], cfg)
}

func boundaryProbeStart(endIdx int, cfg capsuleSegmenterConfig) int {
	window := cfg.OrderingProbeWindow
	if cfg.TopicShiftMinMessages > window {
		window = cfg.TopicShiftMinMessages
	}
	if cfg.TailAuthorWindow > window {
		window = cfg.TailAuthorWindow
	}
	if window <= 0 {
		return 0
	}
	start := endIdx + 1 - window
	if start < 0 {
		return 0
	}
	return start
}

func segmentBoundaryReason(candidate []sqlite.Message, current, next sqlite.Message, cfg capsuleSegmenterConfig) (string, bool) {
	if !localTimestampOrderingSane(candidate, next, cfg) {
		return "", false
	}

	gap := nonNegativeGap(current.Timestamp, next.Timestamp)
	if gap >= cfg.SilenceGap {
		return fmt.Sprintf("silence_gap:%s", gap.Round(time.Minute)), true
	}

	if len(candidate) >= cfg.TopicShiftMinMessages &&
		significantCategoryShift(current.Category, next.Category) &&
		(gap >= cfg.TopicShiftGap || participantShiftLikely(candidate, next, cfg.TailAuthorWindow)) {
		return fmt.Sprintf("topic_shift:%s->%s", normalizeCategory(current.Category), normalizeCategory(next.Category)), true
	}

	return "", false
}

func localTimestampOrderingSane(candidate []sqlite.Message, next sqlite.Message, cfg capsuleSegmenterConfig) bool {
	window := candidate
	if cfg.OrderingProbeWindow > 0 && len(window) > cfg.OrderingProbeWindow {
		window = window[len(window)-cfg.OrderingProbeWindow:]
	}

	totalTransitions := 0
	negativeTransitions := 0
	prev := window[0]
	for i := 1; i < len(window); i++ {
		totalTransitions++
		if window[i].Timestamp.Before(prev.Timestamp) {
			negativeTransitions++
		}
		prev = window[i]
	}
	totalTransitions++
	if next.Timestamp.Before(prev.Timestamp) {
		negativeTransitions++
	}

	if totalTransitions <= 0 {
		return true
	}
	return float64(negativeTransitions)/float64(totalTransitions) <= cfg.MaxNegativeGapRatio
}

func nonNegativeGap(current, next time.Time) time.Duration {
	if next.Before(current) {
		return 0
	}
	return next.Sub(current)
}

func significantCategoryShift(current, next string) bool {
	current = normalizeCategory(current)
	next = normalizeCategory(next)
	if current == "" || next == "" {
		return false
	}
	if current == "uncategorized" || current == "general" || next == "uncategorized" || next == "general" {
		return false
	}
	return current != next
}

func normalizeCategory(category string) string {
	category = strings.ToLower(strings.TrimSpace(category))
	if category == "" {
		return "uncategorized"
	}
	return category
}

func participantShiftLikely(candidate []sqlite.Message, next sqlite.Message, window int) bool {
	if len(candidate) == 0 {
		return false
	}
	if window <= 0 || window > len(candidate) {
		window = len(candidate)
	}

	tail := candidate[len(candidate)-window:]
	nextAuthor := canonicalAuthorKey(next)
	if nextAuthor == "" {
		return false
	}

	for _, msg := range tail {
		if canonicalAuthorKey(msg) == nextAuthor {
			return false
		}
	}
	return true
}

func canonicalAuthorKey(msg sqlite.Message) string {
	if id := strings.TrimSpace(msg.AuthorID); id != "" {
		return id
	}
	return strings.ToLower(strings.TrimSpace(msg.Author))
}

func nextHourlyRunDelay(now time.Time) time.Duration {
	brt := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	next := brt.Add(1 * time.Hour)
	return next.Sub(now)
}
