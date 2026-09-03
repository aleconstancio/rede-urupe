/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package memory

import (
	"fmt"
	"strings"

	"nucleo-engine/internal/pkg/timeutil"
	"google.golang.org/genai"
)

func capsuleSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
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
		Required: []string{"participants", "main_topic", "mood", "typed_facts", "open_loops"},
	}
}

func buildCapsulePrompt(segment capsuleSegment) string {
	var b strings.Builder
	b.WriteString("Resuma estas mensagens como uma capsula de memoria estruturada.\n")
	b.WriteString("O segmento ja foi delimitado por regras deterministicas; nao tente juntar com contexto externo nem inventar continuidade.\n\n")
	b.WriteString("<SEGMENT_INFO>\n")
	b.WriteString(fmt.Sprintf("row_span: %d-%d\n", segment.StartExclusiveRowID+1, segment.EndInclusiveRowID))
	b.WriteString(fmt.Sprintf("message_count: %d\n", segment.MessageCount))
	b.WriteString(fmt.Sprintf("boundary_reason: %s\n", segment.Reason))
	b.WriteString("</SEGMENT_INFO>\n\n")
	b.WriteString("Contrato:\n")
	b.WriteString("- Foque em entendimento duravel e preferencias, nao em detalhe granular de conversa.\n")
	b.WriteString("- Ignore troca rotineira, repeticoes, concordancias obvias e micro-passos sem valor futuro.\n")
	b.WriteString("- typed_facts.events: apenas 0-3 fatos concretos que ainda importariam amanha.\n")
	b.WriteString("- typed_facts.traits: apenas 0-2 preferencias ou padroes realmente estaveis observados nas pessoas.\n")
	b.WriteString("- typed_facts.tensions: apenas 0-2 desacordos, friccoes ou contradicoes ainda vivas.\n")
	b.WriteString("- typed_facts.callbacks: no maximo 1 referencia recorrente realmente memoravel.\n")
	b.WriteString("- open_loops: apenas perguntas, tarefas ou assuntos que ainda merecem follow-up.\n")
	b.WriteString("- Para cada open loop, preencha label, status (open/watch/blocked/resolved), owner e next_step. No maximo 3 loops.\n")
	b.WriteString("- Use apenas participantes realmente presentes nas mensagens.\n")
	b.WriteString("- Se o segmento estiver misto ou disperso, prefira um topico mais honesto e generico; nao force falsa unidade.\n")
	b.WriteString("- Se o bloco tiver pouco valor memoravel, retorne poucos itens em vez de inventar importancia.\n")
	b.WriteString("- Se nao houver open loops reais, retorne lista vazia.\n\n")

	for _, m := range segment.Messages {
		meta := []string{
			fmt.Sprintf("id=%d", m.ID),
			fmt.Sprintf("ts=%s", timeutil.InBrasilia(m.Timestamp).Format("15:04")),
			fmt.Sprintf("cat=%s", m.Category),
		}
		if m.ReplyToID != "" {
			meta = append(meta, "reply")
		}
		if m.IsBot {
			meta = append(meta, "bot")
		}
		b.WriteString(fmt.Sprintf("[%s] %s: %s\n", strings.Join(meta, " "), m.Author, m.Content))
	}
	return b.String()
}
