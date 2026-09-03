/*
 * Copyright (c) 2026 Talos V2 Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	minotaur "github.com/aleconstancio/talos/v2/behavior/persona"
	"nucleo-engine/internal/pkg/timeutil"

	"github.com/google/uuid"
)

func (r *Repository) ensureV33PersonaSeed() error {
	now := timeutil.Now()
	core := minotaur.DefaultCoreIdentityProfile()
	core.IsDefault = false
	coreValuesJSON, err := encodeJSONValue(core.CoreValues)
	if err != nil {
		return err
	}
	if _, err := r.db.Exec(`
		INSERT OR IGNORE INTO core_identity_profiles (
			id, name, display_name, avatar_url, description, identity_prompt, core_values_json, is_enabled, is_default, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, core.ID, core.Name, core.DisplayName, core.AvatarURL, core.Description, core.IdentityPrompt, coreValuesJSON, core.IsEnabled, core.IsDefault, now, now); err != nil {
		return err
	}

	// Seed ERIS
	eris := minotaur.ErisCoreIdentityProfile()
	erisValuesJSON, _ := encodeJSONValue(eris.CoreValues)
	if _, err := r.db.Exec(`
		INSERT OR IGNORE INTO core_identity_profiles (
			id, name, display_name, avatar_url, description, identity_prompt, core_values_json, is_enabled, is_default, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, eris.ID, eris.Name, eris.DisplayName, eris.AvatarURL, eris.Description, eris.IdentityPrompt, erisValuesJSON, eris.IsEnabled, eris.IsDefault, now, now); err != nil {
		return err
	}

	// Seed MICÉLIA 🍄 (Padrão do Núcleo Urupê)
	micelia := minotaur.MiceliaCoreIdentityProfile()
	micelia.IsDefault = true
	miceliaValuesJSON, _ := encodeJSONValue(micelia.CoreValues)
	if _, err := r.db.Exec(`
		INSERT OR IGNORE INTO core_identity_profiles (
			id, name, display_name, avatar_url, description, identity_prompt, core_values_json, is_enabled, is_default, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, micelia.ID, micelia.Name, micelia.DisplayName, micelia.AvatarURL, micelia.Description, micelia.IdentityPrompt, miceliaValuesJSON, micelia.IsEnabled, micelia.IsDefault, now, now); err != nil {
		return err
	}

	for _, overlay := range minotaur.DefaultPersonaOverlays() {
		traitsJSON, err := encodeJSONValue(overlay.Traits)
		if err != nil {
			return err
		}
		intentsJSON, err := encodeJSONValue(overlay.AllowedIntents)
		if err != nil {
			return err
		}
		if _, err := r.db.Exec(`
			INSERT OR IGNORE INTO persona_overlays (
				id, identity_id, name, description, style_prompt, traits_json, allowed_intents_json, is_enabled, is_default, sort_order, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, overlay.ID, overlay.IdentityID, overlay.Name, overlay.Description, overlay.StylePrompt, traitsJSON, intentsJSON, overlay.IsEnabled, overlay.IsDefault, overlay.SortOrder, now, now); err != nil {
			return err
		}
	}
	return nil
}

// ListCoreIdentityProfiles returns all editable bot identities.
func (r *Repository) ListCoreIdentityProfiles() ([]minotaur.CoreIdentityProfile, error) {
	rows, err := r.db.Query(`
		SELECT id, name, display_name, avatar_url, description, identity_prompt, core_values_json, is_enabled, is_default, created_at, updated_at
		FROM core_identity_profiles
		ORDER BY is_default DESC, name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []minotaur.CoreIdentityProfile
	for rows.Next() {
		var item minotaur.CoreIdentityProfile
		var coreValuesJSON string
		if err := rows.Scan(&item.ID, &item.Name, &item.DisplayName, &item.AvatarURL, &item.Description, &item.IdentityPrompt, &coreValuesJSON, &item.IsEnabled, &item.IsDefault, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		item.CoreValues = decodeJSONValue(coreValuesJSON, []string{})
		out = append(out, item)
	}
	return out, rows.Err()
}

// GetCoreIdentityProfile returns one identity by id.
func (r *Repository) GetCoreIdentityProfile(id string) (minotaur.CoreIdentityProfile, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return minotaur.CoreIdentityProfile{}, sql.ErrNoRows
	}

	var item minotaur.CoreIdentityProfile
	var coreValuesJSON string
	err := r.db.QueryRow(`
		SELECT id, name, display_name, avatar_url, description, identity_prompt, core_values_json, is_enabled, is_default, created_at, updated_at
		FROM core_identity_profiles
		WHERE id = ?
	`, id).Scan(&item.ID, &item.Name, &item.DisplayName, &item.AvatarURL, &item.Description, &item.IdentityPrompt, &coreValuesJSON, &item.IsEnabled, &item.IsDefault, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return minotaur.CoreIdentityProfile{}, err
	}
	item.CoreValues = decodeJSONValue(coreValuesJSON, []string{})
	return item, nil
}

// UpsertCoreIdentityProfile stores or updates a core identity profile.
func (r *Repository) UpsertCoreIdentityProfile(profile minotaur.CoreIdentityProfile) error {
	profile.ID = strings.TrimSpace(profile.ID)
	profile.Name = strings.TrimSpace(profile.Name)
	profile.Description = strings.TrimSpace(profile.Description)
	profile.IdentityPrompt = strings.TrimSpace(profile.IdentityPrompt)
	if profile.ID == "" {
		profile.ID = uuid.NewString()
	}
	if profile.Name == "" || profile.IdentityPrompt == "" {
		return errors.New("identity name and prompt are required")
	}

	now := timeutil.Now()
	coreValuesJSON, err := encodeJSONValue(profile.CoreValues)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO core_identity_profiles (
			id, name, display_name, avatar_url, description, identity_prompt, core_values_json, is_enabled, is_default, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			display_name = excluded.display_name,
			avatar_url = excluded.avatar_url,
			description = excluded.description,
			identity_prompt = excluded.identity_prompt,
			core_values_json = excluded.core_values_json,
			is_enabled = excluded.is_enabled,
			is_default = excluded.is_default,
			updated_at = excluded.updated_at
	`, profile.ID, profile.Name, profile.DisplayName, profile.AvatarURL, profile.Description, profile.IdentityPrompt, coreValuesJSON, profile.IsEnabled, profile.IsDefault, now, now)
	return err
}

// ListPersonaOverlays returns all overlays, optionally filtered by identity.
func (r *Repository) ListPersonaOverlays(identityID string) ([]minotaur.PersonaOverlay, error) {
	identityID = strings.TrimSpace(identityID)
	baseQuery := `
		SELECT id, identity_id, name, description, style_prompt, traits_json, allowed_intents_json, is_enabled, is_default, sort_order, created_at, updated_at
		FROM persona_overlays
	`
	var rows *sql.Rows
	var err error
	if identityID == "" {
		rows, err = r.db.Query(baseQuery + ` ORDER BY identity_id ASC, sort_order ASC, name ASC`)
	} else {
		rows, err = r.db.Query(baseQuery+` WHERE identity_id = ? ORDER BY sort_order ASC, name ASC`, identityID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []minotaur.PersonaOverlay
	for rows.Next() {
		var item minotaur.PersonaOverlay
		var traitsJSON, intentsJSON string
		if err := rows.Scan(&item.ID, &item.IdentityID, &item.Name, &item.Description, &item.StylePrompt, &traitsJSON, &intentsJSON, &item.IsEnabled, &item.IsDefault, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		item.Traits = decodeJSONValue(traitsJSON, map[string]float64{})
		item.AllowedIntents = decodeJSONValue(intentsJSON, []string{})
		out = append(out, item)
	}
	return out, rows.Err()
}

// GetPersonaOverlay returns one overlay by id.
func (r *Repository) GetPersonaOverlay(id string) (minotaur.PersonaOverlay, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return minotaur.PersonaOverlay{}, sql.ErrNoRows
	}

	var item minotaur.PersonaOverlay
	var traitsJSON, intentsJSON string
	err := r.db.QueryRow(`
		SELECT id, identity_id, name, description, style_prompt, traits_json, allowed_intents_json, is_enabled, is_default, sort_order, created_at, updated_at
		FROM persona_overlays
		WHERE id = ?
	`, id).Scan(&item.ID, &item.IdentityID, &item.Name, &item.Description, &item.StylePrompt, &traitsJSON, &intentsJSON, &item.IsEnabled, &item.IsDefault, &item.SortOrder, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return minotaur.PersonaOverlay{}, err
	}
	item.Traits = decodeJSONValue(traitsJSON, map[string]float64{})
	item.AllowedIntents = decodeJSONValue(intentsJSON, []string{})
	return item, nil
}

// UpsertPersonaOverlay stores or updates a persona overlay.
func (r *Repository) UpsertPersonaOverlay(overlay minotaur.PersonaOverlay) error {
	overlay.ID = strings.TrimSpace(overlay.ID)
	overlay.IdentityID = strings.TrimSpace(overlay.IdentityID)
	overlay.Name = strings.TrimSpace(overlay.Name)
	overlay.Description = strings.TrimSpace(overlay.Description)
	overlay.StylePrompt = strings.TrimSpace(overlay.StylePrompt)
	if overlay.ID == "" {
		overlay.ID = uuid.NewString()
	}
	if overlay.IdentityID == "" || overlay.Name == "" {
		return errors.New("overlay identity_id and name are required")
	}

	now := timeutil.Now()
	traitsJSON, err := encodeJSONValue(overlay.Traits)
	if err != nil {
		return err
	}
	intentsJSON, err := encodeJSONValue(overlay.AllowedIntents)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO persona_overlays (
			id, identity_id, name, description, style_prompt, traits_json, allowed_intents_json, is_enabled, is_default, sort_order, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			identity_id = excluded.identity_id,
			name = excluded.name,
			description = excluded.description,
			style_prompt = excluded.style_prompt,
			traits_json = excluded.traits_json,
			allowed_intents_json = excluded.allowed_intents_json,
			is_enabled = excluded.is_enabled,
			is_default = excluded.is_default,
			sort_order = excluded.sort_order,
			updated_at = excluded.updated_at
	`, overlay.ID, overlay.IdentityID, overlay.Name, overlay.Description, overlay.StylePrompt, traitsJSON, intentsJSON, overlay.IsEnabled, overlay.IsDefault, overlay.SortOrder, now, now)
	return err
}

// GetPersonaPolicy returns the channel persona policy, seeding a default if necessary.
func (r *Repository) GetPersonaPolicy(channelID string) (minotaur.PersonaPolicy, error) {
	channelID = strings.TrimSpace(channelID)
	if channelID == "" {
		return minotaur.PersonaPolicy{}, sql.ErrNoRows
	}

	if err := r.ensureDefaultPersonaPolicy(channelID); err != nil {
		return minotaur.PersonaPolicy{}, err
	}

	var item minotaur.PersonaPolicy
	var allowedIdentityIDsJSON, allowedPersonaIDsJSON, intentMapJSON, modeMapJSON string
	err := r.db.QueryRow(`
		SELECT channel_id, default_identity_id, default_persona_id, selection_mode,
		       allowed_identity_ids_json, allowed_persona_ids_json, intent_persona_map_json, mode_persona_map_json,
		       manual_override_identity_id, manual_override_persona_id, updated_at
		FROM persona_policy
		WHERE channel_id = ?
	`, channelID).Scan(
		&item.ChannelID,
		&item.DefaultIdentityID,
		&item.DefaultPersonaID,
		&item.SelectionMode,
		&allowedIdentityIDsJSON,
		&allowedPersonaIDsJSON,
		&intentMapJSON,
		&modeMapJSON,
		&item.ManualOverrideIdentityID,
		&item.ManualOverridePersonaID,
		&item.UpdatedAt,
	)
	if err != nil {
		return minotaur.PersonaPolicy{}, err
	}
	item.AllowedIdentityIDs = decodeJSONValue(allowedIdentityIDsJSON, []string{})
	item.AllowedPersonaIDs = decodeJSONValue(allowedPersonaIDsJSON, []string{})
	item.IntentPersonaMap = decodeJSONValue(intentMapJSON, map[string]string{})
	item.ModePersonaMap = decodeJSONValue(modeMapJSON, map[string]string{})
	return item, nil
}

func (r *Repository) ensureDefaultPersonaPolicy(channelID string) error {
	var exists int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM persona_policy WHERE channel_id = ?`, channelID).Scan(&exists); err != nil {
		return err
	}
	if exists > 0 {
		return nil
	}
	return r.UpsertPersonaPolicy(minotaur.DefaultPersonaPolicy(channelID))
}

// UpsertPersonaPolicy stores or updates the persona selection policy.
func (r *Repository) UpsertPersonaPolicy(policy minotaur.PersonaPolicy) error {
	policy.ChannelID = strings.TrimSpace(policy.ChannelID)
	policy.DefaultIdentityID = strings.TrimSpace(policy.DefaultIdentityID)
	policy.DefaultPersonaID = strings.TrimSpace(policy.DefaultPersonaID)
	policy.SelectionMode = strings.TrimSpace(policy.SelectionMode)
	if policy.ChannelID == "" || policy.DefaultIdentityID == "" || policy.DefaultPersonaID == "" {
		return errors.New("channel_id, default_identity_id, and default_persona_id are required")
	}
	if policy.SelectionMode == "" {
		policy.SelectionMode = "fixed"
	}

	now := timeutil.Now()
	allowedIdentityIDsJSON, err := encodeJSONValue(policy.AllowedIdentityIDs)
	if err != nil {
		return err
	}
	allowedPersonaIDsJSON, err := encodeJSONValue(policy.AllowedPersonaIDs)
	if err != nil {
		return err
	}
	intentMapJSON, err := encodeJSONValue(policy.IntentPersonaMap)
	if err != nil {
		return err
	}
	modeMapJSON, err := encodeJSONValue(policy.ModePersonaMap)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO persona_policy (
			channel_id, default_identity_id, default_persona_id, selection_mode,
			allowed_identity_ids_json, allowed_persona_ids_json, intent_persona_map_json, mode_persona_map_json,
			manual_override_identity_id, manual_override_persona_id, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(channel_id) DO UPDATE SET
			default_identity_id = excluded.default_identity_id,
			default_persona_id = excluded.default_persona_id,
			selection_mode = excluded.selection_mode,
			allowed_identity_ids_json = excluded.allowed_identity_ids_json,
			allowed_persona_ids_json = excluded.allowed_persona_ids_json,
			intent_persona_map_json = excluded.intent_persona_map_json,
			mode_persona_map_json = excluded.mode_persona_map_json,
			manual_override_identity_id = excluded.manual_override_identity_id,
			manual_override_persona_id = excluded.manual_override_persona_id,
			updated_at = excluded.updated_at
	`, policy.ChannelID, policy.DefaultIdentityID, policy.DefaultPersonaID, policy.SelectionMode, allowedIdentityIDsJSON, allowedPersonaIDsJSON, intentMapJSON, modeMapJSON, policy.ManualOverrideIdentityID, policy.ManualOverridePersonaID, now)
	return err
}

// ListAdaptivePersonaMemories returns bounded adaptive persona adjustments.
func (r *Repository) ListAdaptivePersonaMemories(channelID string) ([]minotaur.AdaptivePersonaMemory, error) {
	channelID = strings.TrimSpace(channelID)
	rows, err := r.db.Query(`
		SELECT channel_id, identity_id, persona_id, adaptive_style_json, confidence, source, updated_at, expires_at
		FROM adaptive_persona_memory
		WHERE (? = '' OR channel_id = ?)
		ORDER BY updated_at DESC
	`, channelID, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []minotaur.AdaptivePersonaMemory
	for rows.Next() {
		var item minotaur.AdaptivePersonaMemory
		var adaptiveStyleJSON string
		var expiresAt sql.NullTime
		if err := rows.Scan(&item.ChannelID, &item.IdentityID, &item.PersonaID, &adaptiveStyleJSON, &item.Confidence, &item.Source, &item.UpdatedAt, &expiresAt); err != nil {
			return nil, err
		}
		item.AdaptiveStyle = decodeJSONValue(adaptiveStyleJSON, map[string]any{})
		if expiresAt.Valid {
			ts := expiresAt.Time
			item.ExpiresAt = &ts
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

// UpsertAdaptivePersonaMemory stores a bounded adaptive adjustment.
func (r *Repository) UpsertAdaptivePersonaMemory(item minotaur.AdaptivePersonaMemory) error {
	item.ChannelID = strings.TrimSpace(item.ChannelID)
	item.IdentityID = strings.TrimSpace(item.IdentityID)
	item.PersonaID = strings.TrimSpace(item.PersonaID)
	item.Source = strings.TrimSpace(item.Source)
	if item.ChannelID == "" || item.IdentityID == "" || item.PersonaID == "" {
		return errors.New("channel_id, identity_id, and persona_id are required")
	}
	if item.Source == "" {
		item.Source = "manual"
	}

	now := timeutil.Now()
	adaptiveStyleJSON, err := encodeJSONValue(item.AdaptiveStyle)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO adaptive_persona_memory (
			channel_id, identity_id, persona_id, adaptive_style_json, confidence, source, updated_at, expires_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(channel_id, identity_id, persona_id) DO UPDATE SET
			adaptive_style_json = excluded.adaptive_style_json,
			confidence = excluded.confidence,
			source = excluded.source,
			updated_at = excluded.updated_at,
			expires_at = excluded.expires_at
	`, item.ChannelID, item.IdentityID, item.PersonaID, adaptiveStyleJSON, item.Confidence, item.Source, now, item.ExpiresAt)
	return err
}

// ResetAdaptivePersonaMemories clears adaptive memory entries by optional filters.
func (r *Repository) ResetAdaptivePersonaMemories(channelID, identityID, personaID string) error {
	query := `DELETE FROM adaptive_persona_memory WHERE 1=1`
	var args []any
	if channelID = strings.TrimSpace(channelID); channelID != "" {
		query += ` AND channel_id = ?`
		args = append(args, channelID)
	}
	if identityID = strings.TrimSpace(identityID); identityID != "" {
		query += ` AND identity_id = ?`
		args = append(args, identityID)
	}
	if personaID = strings.TrimSpace(personaID); personaID != "" {
		query += ` AND persona_id = ?`
		args = append(args, personaID)
	}
	_, err := r.db.Exec(query, args...)
	return err
}

// UpsertPersonaDeltaProposal stores or updates a persona adaptation proposal.
func (r *Repository) UpsertPersonaDeltaProposal(tx *sql.Tx, prop minotaur.PersonaDeltaProposal) error {
	exec := r.getExecutor(tx)
	prop.ID = strings.TrimSpace(prop.ID)
	prop.ChannelID = strings.TrimSpace(prop.ChannelID)
	prop.IdentityID = strings.TrimSpace(prop.IdentityID)
	prop.Status = strings.TrimSpace(prop.Status)
	if prop.ID == "" {
		prop.ID = uuid.NewString()
	}
	if prop.ChannelID == "" || prop.IdentityID == "" {
		return errors.New("channel_id and identity_id are required")
	}
	if prop.Status == "" {
		prop.Status = "pending"
	}

	now := timeutil.Now()
	changesJSON, err := encodeJSONValue(prop.ProposedChanges)
	if err != nil {
		return err
	}
	evidenceJSON, err := encodeJSONValue(prop.EvidenceMessageIDs)
	if err != nil {
		return err
	}

	_, err = exec.Exec(`
		INSERT INTO persona_delta_proposals (
			id, channel_id, identity_id, persona_id, target, proposed_changes_json, reason,
			evidence_message_ids_json, confidence, status, created_at, updated_at, reviewed_at, expires_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			proposed_changes_json = excluded.proposed_changes_json,
			reason = excluded.reason,
			evidence_message_ids_json = excluded.evidence_message_ids_json,
			confidence = excluded.confidence,
			status = excluded.status,
			updated_at = excluded.updated_at,
			reviewed_at = excluded.reviewed_at,
			expires_at = excluded.expires_at
	`, prop.ID, prop.ChannelID, prop.IdentityID, prop.PersonaID, prop.Target, changesJSON, prop.Reason, evidenceJSON, prop.Confidence, prop.Status, now, now, prop.ReviewedAt, prop.ExpiresAt)
	return err
}

// ListPersonaDeltaProposals returns proposal rows for review.
func (r *Repository) ListPersonaDeltaProposals(channelID, status string) ([]minotaur.PersonaDeltaProposal, error) {
	channelID = strings.TrimSpace(channelID)
	status = strings.TrimSpace(status)
	rows, err := r.db.Query(`
		SELECT id, channel_id, identity_id, persona_id, target, proposed_changes_json, reason,
		       evidence_message_ids_json, confidence, status, created_at, reviewed_at, expires_at
		FROM persona_delta_proposals
		WHERE (? = '' OR channel_id = ?)
		  AND (? = '' OR status = ?)
		ORDER BY created_at DESC
	`, channelID, channelID, status, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []minotaur.PersonaDeltaProposal
	for rows.Next() {
		var item minotaur.PersonaDeltaProposal
		var proposedChangesJSON, evidenceJSON string
		var reviewedAt, expiresAt sql.NullTime
		if err := rows.Scan(&item.ID, &item.ChannelID, &item.IdentityID, &item.PersonaID, &item.Target, &proposedChangesJSON, &item.Reason, &evidenceJSON, &item.Confidence, &item.Status, &item.CreatedAt, &reviewedAt, &expiresAt); err != nil {
			return nil, err
		}
		item.ProposedChanges = decodeJSONValue(proposedChangesJSON, map[string]any{})
		item.EvidenceMessageIDs = decodeJSONValue(evidenceJSON, []string{})
		if reviewedAt.Valid {
			ts := reviewedAt.Time
			item.ReviewedAt = &ts
		}
		if expiresAt.Valid {
			ts := expiresAt.Time
			item.ExpiresAt = &ts
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

// UpdatePersonaDeltaProposalStatus updates review status for a proposal.
func (r *Repository) UpdatePersonaDeltaProposalStatus(id, status string) error {
	id = strings.TrimSpace(id)
	status = strings.TrimSpace(status)
	if id == "" || status == "" {
		return errors.New("proposal id and status are required")
	}
	_, err := r.db.Exec(`
		UPDATE persona_delta_proposals
		SET status = ?, reviewed_at = ?, updated_at = ?
		WHERE id = ?
	`, status, timeutil.Now(), timeutil.Now(), id)
	return err
}

// ApplyPersonaProposal applies an approved proposal to adaptive memory.
func (r *Repository) ApplyPersonaProposal(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("proposal id is required")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var p minotaur.PersonaDeltaProposal
	var changesJSON string
	var existingAdaptiveJSON string
	var existingAdaptive sql.NullString
	err = tx.QueryRow(`
		SELECT channel_id, identity_id, persona_id, proposed_changes_json, confidence, expires_at
		FROM persona_delta_proposals
		WHERE id = ? AND status = 'approved'
	`, id).Scan(&p.ChannelID, &p.IdentityID, &p.PersonaID, &changesJSON, &p.Confidence, &p.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to fetch approved proposal: %w", err)
	}
	p.ProposedChanges = decodeJSONValue(changesJSON, map[string]any{})

	err = tx.QueryRow(`
		SELECT adaptive_style_json
		FROM adaptive_persona_memory
		WHERE channel_id = ? AND identity_id = ? AND persona_id = ?
	`, p.ChannelID, p.IdentityID, p.PersonaID).Scan(&existingAdaptive)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to fetch existing adaptive memory: %w", err)
	}
	if existingAdaptive.Valid {
		existingAdaptiveJSON = existingAdaptive.String
	}

	mergedAdaptive := mergeAdaptiveStyleMaps(
		decodeJSONValue(existingAdaptiveJSON, map[string]any{}),
		p.ProposedChanges,
	)
	mergedAdaptiveJSON, err := encodeJSONValue(mergedAdaptive)
	if err != nil {
		return fmt.Errorf("failed to encode merged adaptive memory: %w", err)
	}

	// Update adaptive memory
	now := timeutil.Now()
	_, err = tx.Exec(`
		INSERT INTO adaptive_persona_memory (
			channel_id, identity_id, persona_id, adaptive_style_json, confidence, source, updated_at, expires_at
		) VALUES (?, ?, ?, ?, ?, 'metacognition', ?, ?)
		ON CONFLICT(channel_id, identity_id, persona_id) DO UPDATE SET
			adaptive_style_json = excluded.adaptive_style_json,
			confidence = excluded.confidence,
			source = excluded.source,
			updated_at = excluded.updated_at,
			expires_at = excluded.expires_at
	`, p.ChannelID, p.IdentityID, p.PersonaID, mergedAdaptiveJSON, p.Confidence, now, p.ExpiresAt)
	if err != nil {
		return fmt.Errorf("failed to update adaptive memory: %w", err)
	}

	// Mark proposal as applied
	_, err = tx.Exec(`UPDATE persona_delta_proposals SET status = 'applied', reviewed_at = ?, updated_at = ? WHERE id = ?`, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to mark proposal as applied: %w", err)
	}

	return tx.Commit()
}

// GetPersonaStudioState builds the operator-facing persona snapshot for a channel.
func (r *Repository) GetPersonaStudioState(channelID string) (minotaur.PersonaStudioState, error) {
	policy, err := r.GetPersonaPolicy(channelID)
	if err != nil {
		return minotaur.PersonaStudioState{}, err
	}

	identityID := policy.DefaultIdentityID
	if strings.TrimSpace(policy.ManualOverrideIdentityID) != "" {
		identityID = strings.TrimSpace(policy.ManualOverrideIdentityID)
	}
	personaID := policy.DefaultPersonaID
	if strings.TrimSpace(policy.ManualOverridePersonaID) != "" {
		personaID = strings.TrimSpace(policy.ManualOverridePersonaID)
	}

	identityProfile, err := r.GetCoreIdentityProfile(identityID)
	if err != nil {
		return minotaur.PersonaStudioState{}, err
	}
	overlay, err := r.GetPersonaOverlay(personaID)
	if err != nil {
		return minotaur.PersonaStudioState{}, err
	}
	adaptiveMemories, err := r.ListAdaptivePersonaMemories(channelID)
	if err != nil {
		return minotaur.PersonaStudioState{}, err
	}

	return minotaur.PersonaStudioState{
		ActiveIdentity:   identityProfile,
		ActiveOverlay:    overlay,
		Policy:           policy,
		AdaptiveMemories: adaptiveMemories,
	}, nil
}

func encodeJSONValue(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func decodeJSONValue[T any](raw string, fallback T) T {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback
	}
	var out T
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return fallback
	}
	return out
}

func mergeAdaptiveStyleMaps(base map[string]any, delta map[string]any) map[string]any {
	if len(base) == 0 && len(delta) == 0 {
		return nil
	}
	merged := make(map[string]any, len(base)+len(delta))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range delta {
		merged[key] = value
	}
	return merged
}
