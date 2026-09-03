/*
 * Copyright (c) 2026 Talos V2 Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"path/filepath"
	"testing"
	"time"

	minotaur "github.com/aleconstancio/talos/v2/behavior/persona"
)

func TestApplyPersonaProposalMergesAdaptiveMemory(t *testing.T) {
	tempDir := t.TempDir()
	repo, err := NewRepository(filepath.Join(tempDir, "talos.db"))
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	defer repo.Close()

	err = repo.UpsertAdaptivePersonaMemory(minotaur.AdaptivePersonaMemory{
		ChannelID:  "channel-1",
		IdentityID: minotaur.DefaultCoreIdentityID,
		PersonaID:  minotaur.DefaultPersonaID,
		AdaptiveStyle: map[string]any{
			"humor_bias":             0.2,
			"preferred_reply_length": "short",
		},
		Confidence: 0.8,
		Source:     "manual",
	})
	if err != nil {
		t.Fatalf("failed to seed adaptive memory: %v", err)
	}

	err = repo.UpsertPersonaDeltaProposal(nil, minotaur.PersonaDeltaProposal{
		ID:         "proposal-1",
		ChannelID:  "channel-1",
		IdentityID: minotaur.DefaultCoreIdentityID,
		PersonaID:  minotaur.DefaultPersonaID,
		Target:     "adaptive",
		ProposedChanges: map[string]any{
			"assertiveness_bias": 0.4,
		},
		Reason:             "better fit",
		EvidenceMessageIDs: []string{"m1"},
		Confidence:         0.9,
		Status:             "approved",
		CreatedAt:          time.Now().UTC(),
	})
	if err != nil {
		t.Fatalf("failed to seed persona proposal: %v", err)
	}

	if err := repo.ApplyPersonaProposal("proposal-1"); err != nil {
		t.Fatalf("failed to apply persona proposal: %v", err)
	}

	memories, err := repo.ListAdaptivePersonaMemories("channel-1")
	if err != nil {
		t.Fatalf("failed to read adaptive memories: %v", err)
	}
	if len(memories) != 1 {
		t.Fatalf("expected 1 adaptive memory row, got %d", len(memories))
	}

	if got := memories[0].AdaptiveStyle["humor_bias"]; got != 0.2 {
		t.Fatalf("expected existing humor_bias to be preserved, got %#v", got)
	}
	if got := memories[0].AdaptiveStyle["preferred_reply_length"]; got != "short" {
		t.Fatalf("expected existing preferred_reply_length to be preserved, got %#v", got)
	}
	if got := memories[0].AdaptiveStyle["assertiveness_bias"]; got != 0.4 {
		t.Fatalf("expected proposal assertiveness_bias to be merged, got %#v", got)
	}

	proposals, err := repo.ListPersonaDeltaProposals("channel-1", "")
	if err != nil {
		t.Fatalf("failed to list proposals: %v", err)
	}
	if len(proposals) != 1 || proposals[0].Status != "applied" {
		t.Fatalf("expected proposal status to become applied, got %+v", proposals)
	}
}

func TestV33SeedsDoNotOverwriteOperatorEdits(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "talos.db")

	repo, err := NewRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	edited := minotaur.DefaultCoreIdentityProfile()
	edited.Name = "Talos Custom"
	edited.IdentityPrompt = "custom prompt"
	if err := repo.UpsertCoreIdentityProfile(edited); err != nil {
		t.Fatalf("failed to write edited identity: %v", err)
	}
	if err := repo.Close(); err != nil {
		t.Fatalf("failed to close repository: %v", err)
	}

	reopened, err := NewRepository(dbPath)
	if err != nil {
		t.Fatalf("failed to reopen repository: %v", err)
	}
	defer reopened.Close()

	got, err := reopened.GetCoreIdentityProfile(minotaur.DefaultCoreIdentityID)
	if err != nil {
		t.Fatalf("failed to read identity after reopen: %v", err)
	}
	if got.Name != "Talos Custom" {
		t.Fatalf("expected custom identity name to survive seed, got %q", got.Name)
	}
	if got.IdentityPrompt != "custom prompt" {
		t.Fatalf("expected custom identity prompt to survive seed, got %q", got.IdentityPrompt)
	}
}
