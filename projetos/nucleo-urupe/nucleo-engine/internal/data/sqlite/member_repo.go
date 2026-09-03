/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"encoding/json"
	"time"
)

type MemberProfile struct {
	DiscordID string   `json:"discord_id"`
	Name      string   `json:"name"`
	Roles     []string `json:"roles"`
	Age       int      `json:"age"`
	Interests string   `json:"interests"`
	Religion  string   `json:"religion"`
	Notes     string   `json:"notes"`
	UpdatedAt string   `json:"updated_at"`
}

func (r *Repository) UpsertMemberProfile(p *MemberProfile) error {
	rolesJSON, _ := json.Marshal(p.Roles)
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := r.db.Exec(`
		INSERT INTO member_profiles (discord_id, name, roles_json, age, interests, religion, notes, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(discord_id) DO UPDATE SET
			name = excluded.name,
			roles_json = excluded.roles_json,
			updated_at = excluded.updated_at
	`, p.DiscordID, p.Name, string(rolesJSON), p.Age, p.Interests, p.Religion, p.Notes, now)
	return err
}

func (r *Repository) GetMemberProfile(discordID string) (*MemberProfile, error) {
	p := &MemberProfile{}
	var rolesJSON string
	err := r.db.QueryRow(`
		SELECT discord_id, name, roles_json, age, interests, religion, notes, updated_at
		FROM member_profiles WHERE discord_id = ?
	`, discordID).Scan(&p.DiscordID, &p.Name, &rolesJSON, &p.Age, &p.Interests, &p.Religion, &p.Notes, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	p.Roles = decodeStringSlice(rolesJSON)
	return p, nil
}

func (r *Repository) GetMemberProfiles(ids []string) ([]MemberProfile, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	query := "SELECT discord_id, name, roles_json, age, interests, religion, notes, updated_at FROM member_profiles WHERE discord_id IN ("
	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		if i > 0 {
			query += ","
		}
		query += "?"
		args = append(args, id)
	}
	query += ")"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []MemberProfile
	for rows.Next() {
		var p MemberProfile
		var rolesJSON string
		if err := rows.Scan(&p.DiscordID, &p.Name, &rolesJSON, &p.Age, &p.Interests, &p.Religion, &p.Notes, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Roles = decodeStringSlice(rolesJSON)
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (r *Repository) ListMemberProfiles(limit, offset int) ([]MemberProfile, error) {
	rows, err := r.db.Query(`
		SELECT discord_id, name, roles_json, age, interests, religion, notes, updated_at
		FROM member_profiles
		ORDER BY updated_at DESC
		LIMIT ? OFFSET ?
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []MemberProfile
	for rows.Next() {
		var p MemberProfile
		var rolesJSON string
		if err := rows.Scan(&p.DiscordID, &p.Name, &rolesJSON, &p.Age, &p.Interests, &p.Religion, &p.Notes, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Roles = decodeStringSlice(rolesJSON)
		profiles = append(profiles, p)
	}
	return profiles, nil
}

func (r *Repository) CountMemberProfiles() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM member_profiles").Scan(&count)
	return count, err
}

func (r *Repository) SearchMemberProfiles(query string) ([]MemberProfile, error) {
	q := "%" + query + "%"
	rows, err := r.db.Query(`
		SELECT discord_id, name, roles_json, age, interests, religion, notes, updated_at
		FROM member_profiles
		WHERE name LIKE ? OR notes LIKE ? OR interests LIKE ?
		ORDER BY updated_at DESC
		LIMIT 50
	`, q, q, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []MemberProfile
	for rows.Next() {
		var p MemberProfile
		var rolesJSON string
		if err := rows.Scan(&p.DiscordID, &p.Name, &rolesJSON, &p.Age, &p.Interests, &p.Religion, &p.Notes, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Roles = decodeStringSlice(rolesJSON)
		profiles = append(profiles, p)
	}
	return profiles, nil
}
