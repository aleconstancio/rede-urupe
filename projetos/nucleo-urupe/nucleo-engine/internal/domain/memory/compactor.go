/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"nucleo-engine/internal/data/llm"
	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/pkg/timeutil"
	"google.golang.org/genai"
)

type Compactor struct {
	repo  *sqlite.Repository
	llm   *llm.Client
	model string
}

type summaryCount struct {
	value string
	count int
}

func NewCompactor(repo *sqlite.Repository, llmClient *llm.Client, model string) *Compactor {
	return &Compactor{
		repo:  repo,
		llm:   llmClient,
		model: model,
	}
}

func (c *Compactor) Start(ctx context.Context) {
	log.Println("[Compactor] Started twice-daily worker")

	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			brt := timeutil.InBrasilia(t)
			hour := brt.Hour()
			min := brt.Minute()
			if (hour == 3 || hour == 15) && min < 30 {
				c.runOnce(ctx)
				time.Sleep(31 * time.Minute)
			}
		}
	}
}

func (c *Compactor) runOnce(ctx context.Context) {
	today := timeutil.Today()
	activeCapsules, err := c.repo.GetActiveMemoryCapsulesByDate(today)
	if err != nil {
		log.Printf("[Compactor] Error fetching capsules: %v", err)
		return
	}

	mergeInput := selectCompactionInput(activeCapsules)
	if len(mergeInput) == 0 {
		log.Printf("[Compactor] Skipping: no eligible capsules for %s", today)
		return
	}

	log.Printf("[Compactor] Merging %d capsules for %s", len(mergeInput), today)

	summaryKind := chooseSummaryKind(mergeInput, strings.ToLower(strings.TrimSpace(c.repo.GetConfig("memory_summary_kind", "auto"))))
	prompt := c.buildPrompt(mergeInput, summaryKind)

	schema := &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"merged_time_span": {Type: genai.TypeString},
			"participants": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
			},
			"main_topic": {Type: genai.TypeString},
			"mood":       {Type: genai.TypeString},
			"typed_facts": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"events": {
						Type:  genai.TypeArray,
						Items: &genai.Schema{Type: genai.TypeString},
					},
					"traits": {
						Type:  genai.TypeArray,
						Items: &genai.Schema{Type: genai.TypeString},
					},
					"tensions": {
						Type:  genai.TypeArray,
						Items: &genai.Schema{Type: genai.TypeString},
					},
					"callbacks": {
						Type:  genai.TypeArray,
						Items: &genai.Schema{Type: genai.TypeString},
					},
				},
				Required: []string{"events", "traits", "tensions", "callbacks"},
			},
			"open_loops": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"label":     {Type: genai.TypeString},
						"status":    {Type: genai.TypeString},
						"owner":     {Type: genai.TypeString},
						"next_step": {Type: genai.TypeString},
					},
					Required: []string{"label", "status", "owner", "next_step"},
				},
			},
		},
		Required: []string{"merged_time_span", "participants", "main_topic", "mood", "typed_facts", "open_loops"},
	}

	completion, err := c.llm.Complete(ctx, []llm.Message{{Role: "user", Content: prompt}}, nil, llm.RequestOptions{
		Model:            c.model,
		ResponseMimeType: "application/json",
		ResponseSchema:   schema,
		Purpose:          "summary",
	})
	if err != nil {
		log.Printf("[Compactor] LLM error: %v", err)
		return
	}

	var merged sqlite.MemoryCapsule
	if err := json.Unmarshal([]byte(completion.Text), &merged); err != nil {
		log.Printf("[Compactor] JSON error: %v", err)
		return
	}

	merged.DayDate = today
	merged.Kind = summaryKind
	merged.IsMerged = false
	merged.Category = dominantCapsuleCategory(mergeInput)
	merged.SourceStartRowID, merged.SourceEndRowID, merged.SourceMessageCount = mergeSourceCoverage(mergeInput)
	if strings.TrimSpace(merged.MainTopic) == "" {
		merged.MainTopic = summaryFocusLabel(mergeInput, summaryKind)
	}
	merged.Normalize()

	tx, err := c.repo.BeginImmediate()
	if err != nil {
		log.Printf("[Compactor] Begin tx error: %v", err)
		return
	}
	defer tx.Rollback()

	if err := c.repo.SaveMemoryCapsuleTx(tx, &merged); err != nil {
		log.Printf("[Compactor] DB error inserting merged capsule: %v", err)
		return
	}

	var ids []int64
	for _, capsule := range mergeInput {
		ids = append(ids, capsule.ID)
	}
	if err := c.repo.MarkCapsulesAsMergedTx(tx, ids); err != nil {
		log.Printf("[Compactor] DB error marking capsules merged: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[Compactor] Commit error: %v", err)
		return
	}

	log.Printf("[Compactor] Saved merged capsule ID %d", merged.ID)
	_ = c.repo.LogBudgetEvent(
		fmt.Sprintf("%s:%s", today, merged.TimeSpan),
		"capsule_compaction",
		c.model,
		"memory",
		completion.Usage.PromptTokens,
		completion.Usage.CompletionTokens,
	)
}

func (c *Compactor) buildPrompt(capsules []sqlite.MemoryCapsule, summaryKind string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Mescle estas capsulas episodicas em uma unica capsula consolidada do tipo %s.\n\n", summaryKind))
	b.WriteString("Instrucoes criticas:\n")
	b.WriteString("1. Nunca invente participantes. Use apenas nomes presentes nas capsulas de origem.\n")
	b.WriteString("2. Preserve tensoes, contradicoes e open loops. Nao resolva conflitos no resumo.\n")
	b.WriteString("3. Mantenha typed_facts organizado em events, traits, tensions e callbacks.\n")
	b.WriteString("4. merged_time_span deve cobrir do inicio da capsula mais antiga ao fim da mais recente.\n")
	b.WriteString("5. Se uma capsula de origem for capture_slice, trate-a como material bruto de baixa prioridade: extraia apenas o que for duravel ou recorrente.\n")
	b.WriteString("6. Open loops resolvidos podem permanecer apenas se ainda forem uteis como contexto; loops vivos devem continuar claros.\n")
	b.WriteString("7. Se o tipo for topic_summary, concentre o resumo no topico mais recorrente e corte ruído lateral.\n")
	b.WriteString("8. Se o tipo for category_summary, priorize o padrao de comportamento da categoria e nao um unico micro-episodio.\n\n")

	for _, cap := range capsules {
		cap.Normalize()
		b.WriteString(fmt.Sprintf("--- CAPSULE [%s | %s | episode=%s] ---\n", cap.TimeSpan, cap.Kind, cap.EpisodeID))
		b.WriteString(fmt.Sprintf("Topic: %s | Mood: %s | Category: %s\n", cap.MainTopic, cap.Mood, cap.Category))
		b.WriteString(fmt.Sprintf("Participants: %s\n", strings.Join(cap.Participants, ", ")))
		b.WriteString(fmt.Sprintf("Typed Facts: %s\n", strings.Join(cap.KeyFacts, "; ")))
		b.WriteString(fmt.Sprintf("Open Loops: %s\n\n", strings.Join(cap.UnresolvedQuestions, "; ")))
	}
	return b.String()
}

func dominantCapsuleCategory(capsules []sqlite.MemoryCapsule) string {
	counts := make(map[string]int)
	order := make([]string, 0, len(capsules))
	for _, cap := range capsules {
		category := cap.Category
		if category == "" {
			category = "general"
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

func selectCompactionInput(active []sqlite.MemoryCapsule) []sqlite.MemoryCapsule {
	var latestSummary *sqlite.MemoryCapsule
	snapshots := make([]sqlite.MemoryCapsule, 0, len(active))

	for i := range active {
		capsule := active[i]
		capsule.Normalize()

		if sqlite.IsSummaryKind(capsule.Kind) {
			if latestSummary == nil || capsule.CreatedAt.After(latestSummary.CreatedAt) {
				copy := capsule
				latestSummary = &copy
			}
			continue
		}
		snapshots = append(snapshots, capsule)
	}

	switch {
	case latestSummary == nil && len(snapshots) <= 3:
		return nil
	case latestSummary != nil && len(snapshots) == 0:
		return nil
	}

	mergeInput := make([]sqlite.MemoryCapsule, 0, len(snapshots)+1)
	if latestSummary != nil {
		mergeInput = append(mergeInput, *latestSummary)
	}
	mergeInput = append(mergeInput, snapshots...)
	return mergeInput
}

func chooseSummaryKind(capsules []sqlite.MemoryCapsule, override string) string {
	switch override {
	case sqlite.MemoryKindDailySummary, sqlite.MemoryKindTopicSummary, sqlite.MemoryKindCategorySummary:
		return override
	}

	if kind := chooseSummaryKindHeuristic(capsules); kind != "" {
		return kind
	}
	return sqlite.MemoryKindDailySummary
}

func chooseSummaryKindHeuristic(capsules []sqlite.MemoryCapsule) string {
	if len(capsules) == 0 {
		return ""
	}

	categoryCounts := make(map[string]int)
	topicCounts := make(map[string]int)
	totalSnapshots := 0
	for _, capsule := range capsules {
		capsule.Normalize()
		if sqlite.IsSummaryKind(capsule.Kind) {
			continue
		}
		totalSnapshots++

		if category := normalizeSummaryValue(capsule.Category); category != "" && category != "general" && category != "uncategorized" {
			categoryCounts[category]++
		}
		if topic := normalizeSummaryValue(capsule.MainTopic); topic != "" && topic != "mixed conversation segment" {
			topicCounts[topic]++
		}
	}

	if totalSnapshots == 0 {
		return ""
	}

	bestCategory := maxCountValue(categoryCounts)
	if bestCategory.count >= 2 && float64(bestCategory.count)/float64(totalSnapshots) >= 0.60 {
		return sqlite.MemoryKindCategorySummary
	}

	bestTopic := maxCountValue(topicCounts)
	if bestTopic.count >= 2 && float64(bestTopic.count)/float64(totalSnapshots) >= 0.50 {
		return sqlite.MemoryKindTopicSummary
	}

	return sqlite.MemoryKindDailySummary
}

func normalizeSummaryValue(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.Trim(v, "[](){}<>")
	v = strings.Join(strings.Fields(v), " ")
	return v
}

func maxCountValue(counts map[string]int) summaryCount {
	best := summaryCount{}
	for value, count := range counts {
		if count > best.count || (count == best.count && value < best.value) {
			best = summaryCount{value: value, count: count}
		}
	}
	return best
}

func summaryFocusLabel(capsules []sqlite.MemoryCapsule, kind string) string {
	switch kind {
	case sqlite.MemoryKindCategorySummary:
		return dominantCapsuleCategory(capsules)
	case sqlite.MemoryKindTopicSummary:
		return dominantCapsuleTopic(capsules)
	default:
		return dominantCapsuleCategory(capsules)
	}
}

func mergeSourceCoverage(capsules []sqlite.MemoryCapsule) (int64, int64, int) {
	var minStart int64
	var maxEnd int64
	totalMessages := 0

	for _, capsule := range capsules {
		if minStart == 0 || (capsule.SourceStartRowID > 0 && capsule.SourceStartRowID < minStart) {
			minStart = capsule.SourceStartRowID
		}
		if capsule.SourceEndRowID > maxEnd {
			maxEnd = capsule.SourceEndRowID
		}
		totalMessages += capsule.SourceMessageCount
	}

	return minStart, maxEnd, totalMessages
}

func dominantCapsuleTopic(capsules []sqlite.MemoryCapsule) string {
	counts := make(map[string]int)
	order := make([]string, 0, len(capsules))
	for _, cap := range capsules {
		topic := normalizeSummaryValue(cap.MainTopic)
		if topic == "" || topic == "mixed conversation segment" {
			continue
		}
		if _, seen := counts[topic]; !seen {
			order = append(order, topic)
		}
		counts[topic]++
	}

	best := ""
	bestCount := -1
	for _, topic := range order {
		if counts[topic] > bestCount {
			best = topic
			bestCount = counts[topic]
		}
	}
	return best
}
