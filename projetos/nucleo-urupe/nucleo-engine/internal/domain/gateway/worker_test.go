package gateway

import (
	"strings"
	"testing"
	"time"

	"nucleo-engine/internal/data/sqlite"
	minotaur "github.com/aleconstancio/talos/v2/behavior/persona"
	"nucleo-engine/internal/pkg/timeutil"
)

func TestExtractKeywordsFiltersNoiseAndCapsCount(t *testing.T) {
	messages := []sqlite.Message{
		{Content: "the architecture around sqlite worker queue http://example.com 123"},
		{Content: "working memory prefetch extracts keywords from recent messages"},
	}

	keywords := extractKeywords(messages)
	if len(keywords) == 0 {
		t.Fatal("expected keywords to be extracted")
	}
	if len(keywords) > 6 {
		t.Fatalf("expected at most 6 keywords, got %d: %v", len(keywords), keywords)
	}

	want := map[string]bool{
		"architecture": true,
		"sqlite":       true,
		"worker":       true,
		"queue":        true,
	}
	for _, kw := range keywords {
		if kw == "the" || kw == "http" || kw == "123" {
			t.Fatalf("unexpected stopword or noise keyword in %v", keywords)
		}
		delete(want, kw)
	}
	if len(want) != 0 {
		t.Fatalf("missing expected keywords: %v from %v", want, keywords)
	}
}

func TestBuildGateSystemPromptEnforcesSilenceFirstTriage(t *testing.T) {
	a := &PayloadAssembler{}
	prompt := a.BuildGateSystemPrompt("Talos")

	for _, snippet := range []string{
		"Silence is normal",
		"closed_human_loop",
		"rhetorical_question",
		"insider_joke",
		"False positives are socially expensive",
	} {
		if !strings.Contains(prompt, snippet) {
			t.Fatalf("expected gate prompt to contain %q", snippet)
		}
	}
}

func TestAssembleReplyPromptIncludesVibeMatchingAndMonologueContract(t *testing.T) {
	a := &PayloadAssembler{}
	resolved := minotaur.ResolvedPersona{
		Identity: minotaur.DefaultCoreIdentityProfile(),
		Overlays: []minotaur.PersonaOverlay{minotaur.DefaultPersonaOverlays()[0]},
	}
	prompt := a.AssembleReplyPrompt(resolved, false)

	for _, snippet := range []string{
		"VIBE MATCHING",
		"Match the room's message length",
		"internal_monologue",
		"frame_control",
		"interrupt_check",
		"grounding_ledger",
	} {
		if !strings.Contains(prompt, snippet) {
			t.Fatalf("expected reply prompt to contain %q", snippet)
		}
	}
}

func TestBuildReplyUserPromptIncludesGateDiagnosis(t *testing.T) {
	a := &PayloadAssembler{}
	userPrompt := a.BuildReplyUserPrompt(PromptContext{}, false, &GateResponse{
		ReasonCode:       "additive_contrast",
		Confidence:       0.88,
		SocialFrame:      "open_floor",
		AdditiveValue:    "high",
		InterruptionRisk: "low",
	}, timeutil.InBrasilia(time.Unix(0, 0)))

	for _, snippet := range []string{
		"<GATE_DIAGNOSIS>",
		"reason: additive_contrast",
		"social_frame: open_floor",
		"additive_value: high",
		"interruption_risk: low",
	} {
		if !strings.Contains(userPrompt, snippet) {
			t.Fatalf("expected reply user prompt to contain %q", snippet)
		}
	}
}

func TestWriteThreeLayerContextIncludesActiveOpenLoops(t *testing.T) {
	var b strings.Builder
	writeThreeLayerContext(&b, PromptContext{
		Episodic: []sqlite.MemoryCapsule{
			{
				ID:        7,
				DayDate:   "2026-05-08",
				TimeSpan:  "12:00-12:10",
				EpisodeID: "2026-05-08:episode_snapshot:1-4",
				Kind:      sqlite.MemoryKindEpisodeSnapshot,
				MainTopic: "memory retrieval",
				OpenLoops: []sqlite.OpenLoop{
					{Label: "validate open loop ranking", Status: "open", Owner: "Talos", NextStep: "add tests"},
				},
			},
		},
		OpenLoops: []sqlite.OpenLoop{
			{Label: "validate open loop ranking", Status: "open", Owner: "Talos", NextStep: "add tests"},
		},
		Working: []sqlite.Message{
			{ID: 42, Timestamp: timeutil.InBrasilia(time.Unix(0, 0)), Author: "Ale", Content: "retoma a parte do ranking", EpisodeID: "2026-05-08:episode_snapshot:1-4"},
		},
	}, timeutil.InBrasilia(time.Unix(0, 0)))

	out := b.String()
	for _, snippet := range []string{
		"<ACTIVE_OPEN_LOOPS",
		"validate open loop ranking",
		"owner=Talos",
		"episode=2026-05-08:episode_snapshot:1-4",
	} {
		if !strings.Contains(out, snippet) {
			t.Fatalf("expected context to contain %q, got:\n%s", snippet, out)
		}
	}
}

func TestSelectEpisodicCapsulesPrefersSummaryAndRelevantSnapshots(t *testing.T) {
	capsules := []sqlite.MemoryCapsule{
		{
			ID:        1,
			DayDate:   "2026-05-08",
			TimeSpan:  "09:00-09:30",
			EpisodeID: "2026-05-08:daily_summary:1-30",
			Kind:      sqlite.MemoryKindDailySummary,
			MainTopic: "day summary",
		},
		{
			ID:             11,
			DayDate:        "2026-05-08",
			TimeSpan:       "10:00-10:20",
			EpisodeID:      "2026-05-08:topic_summary:31-42",
			Kind:           sqlite.MemoryKindTopicSummary,
			MainTopic:      "episodic memory retrieval",
			SourceEndRowID: 42,
		},
		{
			ID:             12,
			DayDate:        "2026-05-08",
			TimeSpan:       "10:25-10:50",
			EpisodeID:      "2026-05-08:category_summary:43-55",
			Kind:           sqlite.MemoryKindCategorySummary,
			MainTopic:      "memory routing",
			Category:       "memory",
			SourceEndRowID: 55,
		},
		{
			ID:                  2,
			DayDate:             "2026-05-08",
			TimeSpan:            "12:00-12:20",
			EpisodeID:           "2026-05-08:episode_snapshot:31-50",
			Kind:                sqlite.MemoryKindEpisodeSnapshot,
			MainTopic:           "episodic memory retrieval",
			KeyFacts:            []string{"EVENT|retrieval ranking changed"},
			UnresolvedQuestions: []string{"OPEN_LOOP|validate ranking"},
			SourceEndRowID:      50,
		},
		{
			ID:             3,
			DayDate:        "2026-05-08",
			TimeSpan:       "13:00-13:10",
			EpisodeID:      "2026-05-08:episode_snapshot:51-60",
			Kind:           sqlite.MemoryKindEpisodeSnapshot,
			MainTopic:      "banter",
			SourceEndRowID: 60,
		},
		{
			ID:             4,
			DayDate:        "2026-05-08",
			TimeSpan:       "13:10-13:40",
			EpisodeID:      "2026-05-08:capture_slice:61-95",
			Kind:           sqlite.MemoryKindCaptureSlice,
			MainTopic:      "thin chatter",
			SourceEndRowID: 95,
		},
	}

	selected := selectEpisodicCapsules(capsules, []string{"retrieval", "ranking"}, 4)
	if len(selected) != 4 {
		t.Fatalf("expected 4 capsules, got %d", len(selected))
	}
	if selected[0].Kind != sqlite.MemoryKindDailySummary {
		t.Fatalf("expected summary to be retained as baseline, got %+v", selected)
	}

	foundRelevant := false
	for _, capsule := range selected {
		if capsule.ID == 2 {
			foundRelevant = true
		}
	}
	if !foundRelevant {
		t.Fatalf("expected relevant snapshot to be selected, got %+v", selected)
	}
	foundTopicSummary := false
	foundCategorySummary := false
	for _, capsule := range selected {
		if capsule.Kind == sqlite.MemoryKindTopicSummary {
			foundTopicSummary = true
		}
		if capsule.Kind == sqlite.MemoryKindCategorySummary {
			foundCategorySummary = true
		}
	}
	if !foundTopicSummary || !foundCategorySummary {
		t.Fatalf("expected topic and category summaries to be selectable, got %+v", selected)
	}
	for _, capsule := range selected {
		if capsule.Kind == sqlite.MemoryKindCaptureSlice {
			t.Fatalf("did not expect capture slices in episodic retrieval, got %+v", selected)
		}
	}
}

func TestReplySchemaRequiresInternalMonologue(t *testing.T) {
	schema := ReplySchema()
	internal, ok := schema.Properties["internal_monologue"]
	if !ok {
		t.Fatal("expected internal_monologue property in reply schema")
	}
	if internal.Type == "" {
		t.Fatal("expected internal_monologue schema type to be populated")
	}
	for _, key := range []string{"surface_read", "subtext", "frame_control", "interrupt_check", "value_add", "memory_decision", "vibe_plan", "draft_plan"} {
		if _, ok := internal.Properties[key]; !ok {
			t.Fatalf("expected internal_monologue to require %q", key)
		}
	}
}

func TestParseReplyResponseSupportsStructuredMonologue(t *testing.T) {
	raw := `{
		"internal_monologue": {
			"surface_read": "direct question about architecture",
			"subtext": "they want a clean answer, not a speech",
			"frame_control": "none",
			"interrupt_check": "safe to join",
			"value_add": "give the missing distinction",
			"memory_decision": "working only",
			"vibe_plan": "2 short sentences, neutral slang, low humor",
			"draft_plan": "answer directly and stop"
		},
		"grounding_ledger": ["working:id=12 direct question"],
		"should_intervene": true,
		"reply_text": "E a diferenca entre fila e worker unico.",
		"stance_updates": []
	}`

	resp, err := ParseReplyResponse(raw)
	if err != nil {
		t.Fatalf("expected structured reply to parse, got error: %v", err)
	}
	if resp.InternalMonologue.InterruptCheck != "safe to join" {
		t.Fatalf("unexpected interrupt_check: %+v", resp.InternalMonologue)
	}
	if len(resp.GroundingLedger) != 1 || resp.GroundingLedger[0] != "working:id=12 direct question" {
		t.Fatalf("unexpected grounding ledger: %+v", resp.GroundingLedger)
	}
}
