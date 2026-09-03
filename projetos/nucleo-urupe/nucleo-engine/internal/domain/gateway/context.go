/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/pkg/timeutil"
)

type PromptContext struct {
	Associative  []sqlite.MemoryCapsule
	Episodic     []sqlite.MemoryCapsule
	OpenLoops    []sqlite.OpenLoop
	Working      []sqlite.Message
	Participants []sqlite.MemberProfile
}

func (a *PayloadAssembler) BuildPromptContext(channelID string, working []sqlite.Message, now time.Time) (PromptContext, error) {
	today := timeutil.InBrasilia(now).Format("2006-01-02")

	episodicAll, err := a.repo.GetActiveMemoryCapsulesByDate(today)
	if err != nil {
		return PromptContext{}, fmt.Errorf("load episodic memory: %w", err)
	}

	prefetch := working
	if len(prefetch) > 5 {
		prefetch = prefetch[len(prefetch)-5:]
	}
	keywords := extractKeywords(prefetch)
	episodic := selectEpisodicCapsules(episodicAll, keywords, 4)
	openLoops := collectOpenLoops(episodic, 5)

	associative, err := a.repo.SearchPastCapsules(keywords, today, 2)
	if err != nil {
		return PromptContext{}, fmt.Errorf("load associative memory: %w", err)
	}

	// Resolve unique participants from working memory
	participantIDs := make(map[string]struct{})
	for _, m := range working {
		if m.AuthorID != "" && !m.IsBot {
			participantIDs[m.AuthorID] = struct{}{}
		}
	}
	var ids []string
	for id := range participantIDs {
		ids = append(ids, id)
	}
	participants, _ := a.repo.GetMemberProfiles(ids)

	return PromptContext{
		Associative:  associative,
		Episodic:     episodic,
		OpenLoops:    openLoops,
		Working:      working,
		Participants: participants,
	}, nil
}

func (a *PayloadAssembler) BuildReplyUserPrompt(ctx PromptContext, reactive bool, gateResp *GateResponse, now time.Time) string {
	var b strings.Builder
	if reactive {
		b.WriteString("The bot was directly addressed.\n\n")
	} else {
		b.WriteString("The ambient gate approved intervention.\n")
		if gateResp != nil {
			b.WriteString(fmt.Sprintf("<GATE_DIAGNOSIS>\nreason: %s\nconfidence: %.2f\nsocial_frame: %s\nadditive_value: %s\ninterruption_risk: %s\n</GATE_DIAGNOSIS>\n\n", gateResp.ReasonCode, gateResp.Confidence, gateResp.SocialFrame, gateResp.AdditiveValue, gateResp.InterruptionRisk))
		}
	}

	b.WriteString("1. Prioritize WORKING_MEMORY for the active dialogue flow.\n")
	b.WriteString("2. Use EPISODIC_MEMORY to catch up on the day's context or topics mentioned slightly earlier.\n")
	b.WriteString("3. ACTIVE_OPEN_LOOPS are follow-up hints, not a to-do list. Mention them only when the room naturally reopens them.\n")
	b.WriteString("4. If a topic in WORKING_MEMORY refers to something in EPISODIC_MEMORY, bridge them naturally: 'as we were saying earlier...', 'regarding what came up just now...'.\n")
	b.WriteString("5. ASSOCIATIVE_MEMORY is for historical analogies only. Never treat it as part of today's conversation.\n\n")

	writeThreeLayerContext(&b, ctx, now)
	return b.String()
}

func (a *PayloadAssembler) BuildGateUserPrompt(ctx PromptContext, now time.Time) string {
	var b strings.Builder
	b.WriteString("Review the live room state below.\n")
	b.WriteString("If the bot is not clearly invited and does not have distinct additive value, remain silent.\n\n")
	b.WriteString("ACTIVE_OPEN_LOOPS are only weak follow-up hints. They do not override the silence-first gate.\n\n")
	writeThreeLayerContext(&b, ctx, now)
	return b.String()
}

func (a *PayloadAssembler) BuildUserPromptPayload(ctx context.Context, messageID string, basePrompt string) interface{} {
	// Vision disabled to save tokens. Only filenames in text context are provided.
	return basePrompt
}

func writeThreeLayerContext(b *strings.Builder, ctx PromptContext, now time.Time) {
	b.WriteString(fmt.Sprintf("<RUNTIME now_brt=\"%s\" today_brt=\"%s\">\n\n", timeutil.InBrasilia(now).Format(time.RFC3339), timeutil.InBrasilia(now).Format("2006-01-02")))

	if len(ctx.Participants) > 0 {
		b.WriteString("<PARTICIPANTS usage=\"social_intelligence_only\">\n")
		for _, p := range ctx.Participants {
			var details []string
			if p.Age > 0 {
				details = append(details, fmt.Sprintf("age=%d", p.Age))
			}
			if p.Interests != "" {
				details = append(details, fmt.Sprintf("interests=%s", p.Interests))
			}
			if p.Religion != "" {
				details = append(details, fmt.Sprintf("religion=%s", p.Religion))
			}
			if len(p.Roles) > 0 {
				details = append(details, fmt.Sprintf("roles=[%s]", strings.Join(p.Roles, ",")))
			}
			if p.Notes != "" {
				details = append(details, fmt.Sprintf("notes=%s", p.Notes))
			}

			b.WriteString(fmt.Sprintf("  - %s (%s): %s\n", p.Name, p.DiscordID, strings.Join(details, " | ")))
		}
		b.WriteString("</PARTICIPANTS>\n\n")
	}

	b.WriteString("<ASSOCIATIVE_MEMORY usage=\"historical_analogy_only\">\n")
	if len(ctx.Associative) == 0 {
		b.WriteString("  (none)\n")
	} else {
		for _, c := range ctx.Associative {
			b.WriteString(formatCapsule("Past analog", c))
		}
	}
	b.WriteString("</ASSOCIATIVE_MEMORY>\n\n")

	b.WriteString("<EPISODIC_MEMORY usage=\"same_day_summary_only\">\n")
	if len(ctx.Episodic) == 0 {
		b.WriteString("  (none)\n")
	} else {
		for _, c := range ctx.Episodic {
			b.WriteString(formatCapsule("Today", c))
		}
	}
	b.WriteString("</EPISODIC_MEMORY>\n\n")

	b.WriteString("<ACTIVE_OPEN_LOOPS usage=\"follow_up_hints_only\">\n")
	if len(ctx.OpenLoops) == 0 {
		b.WriteString("  (none)\n")
	} else {
		for _, loop := range ctx.OpenLoops {
			status := strings.TrimSpace(loop.Status)
			if status == "" {
				status = "open"
			}
			owner := strings.TrimSpace(loop.Owner)
			nextStep := strings.TrimSpace(loop.NextStep)
			b.WriteString(fmt.Sprintf("  - [%s] %s", status, loop.Label))
			if owner != "" {
				b.WriteString(fmt.Sprintf(" | owner=%s", owner))
			}
			if nextStep != "" {
				b.WriteString(fmt.Sprintf(" | next=%s", nextStep))
			}
			b.WriteString("\n")
		}
	}
	b.WriteString("</ACTIVE_OPEN_LOOPS>\n\n")

	b.WriteString("<WORKING_MEMORY usage=\"live_dialogue_only\">\n")
	if len(ctx.Working) == 0 {
		b.WriteString("  (none)\n")
	} else {
		for _, m := range ctx.Working {
			author := m.Author
			if m.IsBot {
				author = "BOT"
			}

			var meta strings.Builder
			meta.WriteString(fmt.Sprintf("id=%d, ts=%s", m.ID, timeutil.InBrasilia(m.Timestamp).Format("15:04")))
			if strings.TrimSpace(m.EpisodeID) != "" {
				meta.WriteString(fmt.Sprintf(", episode=%s", m.EpisodeID))
			}

			if m.ReplyToID != "" {
				refID := m.ReplyToID
				for _, prev := range ctx.Working {
					if prev.DiscordID == m.ReplyToID {
						refID = fmt.Sprintf("%d", prev.ID)
						break
					}
				}
				meta.WriteString(fmt.Sprintf(", reply_to=%s", refID))
			}

			if len(m.Reactions) > 0 {
				meta.WriteString(fmt.Sprintf(", reactions=[%s]", strings.Join(m.Reactions, ",")))
			}

			content := m.Content
			if len(m.Attachments) > 0 {
				var files []string
				for _, a := range m.Attachments {
					parts := strings.Split(a, "/")
					fname := parts[len(parts)-1]
					if idx := strings.Index(fname, "?"); idx != -1 {
						fname = fname[:idx]
					}
					files = append(files, fname)
				}
				content = fmt.Sprintf("%s [MEDIA: %s]", content, strings.Join(files, ", "))
			}

			b.WriteString(fmt.Sprintf("  - [%s] %s (%s): %s\n", meta.String(), author, m.Category, content))
		}
	}
	b.WriteString("</WORKING_MEMORY>\n")
}

func formatCapsule(label string, c sqlite.MemoryCapsule) string {
	c.Normalize()
	return fmt.Sprintf(
		"  - [%s, id=%d, kind=%s, episode=%s, span=%s, rows=%d-%d] day=%s topic=%s participants=[%s] facts=[%s] open_loops=[%s]\n",
		label,
		c.ID,
		c.Kind,
		c.EpisodeID,
		c.TimeSpan,
		c.SourceStartRowID,
		c.SourceEndRowID,
		c.DayDate,
		c.MainTopic,
		strings.Join(c.Participants, ", "),
		strings.Join(c.KeyFacts, "; "),
		strings.Join(c.UnresolvedQuestions, "; "),
	)
}
