/*
 * Copyright (c) 2026 Frente Urupê
 * Licensed under the MIT License.
 */

package sqlite

import (
	"database/sql"

	"nucleo-engine/internal/pkg/timeutil"
)

type ManifestoVersion struct {
	ID          string `json:"id"`
	Version     string `json:"version"` // e.g. "v1.0", "v1.1", "v1.2"
	Title       string `json:"title"`
	Changelog   string `json:"changelog"`
	Content     string `json:"content"`
	Author      string `json:"author"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

func (r *Repository) ensureManifestoTables() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS manifesto_versions (
			id TEXT PRIMARY KEY,
			version TEXT UNIQUE NOT NULL,
			title TEXT NOT NULL,
			changelog TEXT,
			content TEXT NOT NULL,
			author TEXT NOT NULL,
			is_active INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_manifesto_version ON manifesto_versions(version);
		CREATE INDEX IF NOT EXISTS idx_manifesto_active ON manifesto_versions(is_active);
	`)
	return err
}

func (r *Repository) ensureManifestoSeed() error {
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM manifesto_versions").Scan(&count)
	if count > 0 {
		return nil
	}

	initial := ManifestoVersion{
		ID:        "man-v1.0",
		Version:   "v1.0",
		Title:     "Manifesto Ecossocialista por uma Soberania Digital e Emancipação Popular",
		Changelog: "Edição fundacional da Frente Urupê.",
		Content:   "A Frente Urupê nasce do imperativo histórico de retomar o controle sobre a nossa comunicação, nosso trabalho e nosso território. Não aceitamos ser meros produtos de algoritmos corporativos desenhados para extrair atenção e gerar ansiedade. A Rede Urupê é a nossa praça pública digital soberana, antifrágil e construída pelo povo para o povo.",
		Author:    "Coordenação Nacional",
		IsActive:  true,
	}

	return r.CreateManifestoVersion(initial)
}

func (r *Repository) ListManifestoVersions() ([]ManifestoVersion, error) {
	rows, err := r.db.Query(`
		SELECT id, version, title, changelog, content, author, is_active, created_at
		FROM manifesto_versions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []ManifestoVersion
	for rows.Next() {
		var v ManifestoVersion
		var active int
		if err := rows.Scan(&v.ID, &v.Version, &v.Title, &v.Changelog, &v.Content, &v.Author, &active, &v.CreatedAt); err != nil {
			return nil, err
		}
		v.IsActive = active == 1
		versions = append(versions, v)
	}

	return versions, nil
}

func (r *Repository) GetActiveManifesto() (*ManifestoVersion, error) {
	row := r.db.QueryRow(`
		SELECT id, version, title, changelog, content, author, is_active, created_at
		FROM manifesto_versions
		WHERE is_active = 1
		LIMIT 1
	`)

	var v ManifestoVersion
	var active int
	if err := row.Scan(&v.ID, &v.Version, &v.Title, &v.Changelog, &v.Content, &v.Author, &active, &v.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	v.IsActive = active == 1

	return &v, nil
}

func (r *Repository) CreateManifestoVersion(v ManifestoVersion) error {
	now := timeutil.Now()
	activeInt := 0
	if v.IsActive {
		activeInt = 1
		// Deactivate all existing versions if this one is active
		_, _ = r.db.Exec(`UPDATE manifesto_versions SET is_active = 0`)
	}

	_, err := r.db.Exec(`
		INSERT INTO manifesto_versions (id, version, title, changelog, content, author, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			version = excluded.version,
			title = excluded.title,
			changelog = excluded.changelog,
			content = excluded.content,
			author = excluded.author,
			is_active = excluded.is_active
	`, v.ID, v.Version, v.Title, v.Changelog, v.Content, v.Author, activeInt, now)

	return err
}

func (r *Repository) ActivateManifestoVersion(versionID string) error {
	_, err := r.db.Exec(`UPDATE manifesto_versions SET is_active = 0`)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`UPDATE manifesto_versions SET is_active = 1 WHERE id = ?`, versionID)
	return err
}
