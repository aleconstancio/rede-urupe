/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"google.golang.org/genai"
)

func GateSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"social_frame":      {Type: genai.TypeString},
			"invitation_signal": {Type: genai.TypeString},
			"additive_value":    {Type: genai.TypeString},
			"interruption_risk": {Type: genai.TypeString},
			"thinking_ledger":   {Type: genai.TypeString},
			"should_intervene":  {Type: genai.TypeBoolean},
			"reason_code":       {Type: genai.TypeString},
			"confidence":        {Type: genai.TypeNumber},
		},
		Required: []string{"social_frame", "invitation_signal", "additive_value", "interruption_risk", "thinking_ledger", "should_intervene", "reason_code", "confidence"},
	}
}

func ReplySchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"internal_monologue": {
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"surface_read":    {Type: genai.TypeString},
					"subtext":         {Type: genai.TypeString},
					"frame_control":   {Type: genai.TypeString},
					"interrupt_check": {Type: genai.TypeString},
					"value_add":       {Type: genai.TypeString},
					"memory_decision": {Type: genai.TypeString},
					"vibe_plan":       {Type: genai.TypeString},
					"draft_plan":      {Type: genai.TypeString},
				},
				Required: []string{"surface_read", "subtext", "frame_control", "interrupt_check", "value_add", "memory_decision", "vibe_plan", "draft_plan"},
			},
			"grounding_ledger": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
			},
			"should_intervene": {Type: genai.TypeBoolean},
			"reply_text":       {Type: genai.TypeString},
			"suggested_reactions": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
			},
			"stance_updates": {
				Type:  genai.TypeArray,
				Items: &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"internal_monologue", "grounding_ledger", "should_intervene", "reply_text", "suggested_reactions", "stance_updates"},
	}
}
