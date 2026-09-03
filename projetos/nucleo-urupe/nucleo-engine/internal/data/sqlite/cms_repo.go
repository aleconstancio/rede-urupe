/*
 * Copyright (c) 2026 Frente Urupê
 * Licensed under the MIT License.
 */

package sqlite

import (
	"database/sql"
	"strings"

	"nucleo-engine/internal/pkg/timeutil"
)

type Article struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	Content     string `json:"content"`
	Author      string `json:"author"`
	Category    string `json:"category"`
	CoverImage  string `json:"cover_image"`
	IsPublished bool   `json:"is_published"`
	P2PSigned   bool   `json:"p2p_signed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (r *Repository) ensureCMSTables() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS cms_articles (
			id TEXT PRIMARY KEY,
			slug TEXT UNIQUE NOT NULL,
			title TEXT NOT NULL,
			summary TEXT,
			content TEXT NOT NULL,
			author TEXT NOT NULL,
			category TEXT NOT NULL,
			cover_image TEXT,
			is_published INTEGER NOT NULL DEFAULT 1,
			p2p_signed INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_cms_articles_slug ON cms_articles(slug);
		CREATE INDEX IF NOT EXISTS idx_cms_articles_published ON cms_articles(is_published);
	`)
	return err
}

func (r *Repository) ensureCMSArticlesSeed() error {
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM cms_articles").Scan(&count)
	if count > 0 {
		return nil
	}

	manifesto := Article{
		ID:          "art-001",
		Slug:        "manifesto-frente-urupe",
		Title:       "Manifesto por uma Soberania Digital Ecossocialista",
		Summary:     "Construindo redes livres, tecnologia descentralizada e emancipação popular sem a exploração das Big Techs.",
		Content:     "A Frente Urupê nasce do imperativo histórico de retomar o controle sobre a nossa comunicação, nosso trabalho e nosso território. Não aceitamos ser meros produtos de algoritmos corporativos desenhados para extrair atenção e gerar ansiedade. A Rede Urupê é a nossa praça pública digital soberana, antifrágil e construída pelo povo para o povo.",
		Author:      "Frente Urupê",
		Category:    "Manifesto",
		CoverImage:  "",
		IsPublished: true,
		P2PSigned:   true,
	}

	return r.UpsertArticle(manifesto)
}

func (r *Repository) ListArticles(publishedOnly bool) ([]Article, error) {
	query := `SELECT id, slug, title, summary, content, author, category, cover_image, is_published, p2p_signed, created_at, updated_at FROM cms_articles`
	if publishedOnly {
		query += ` WHERE is_published = 1`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		var pub, p2p int
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Summary, &a.Content, &a.Author, &a.Category, &a.CoverImage, &pub, &p2p, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		a.IsPublished = pub == 1
		a.P2PSigned = p2p == 1
		articles = append(articles, a)
	}

	return articles, nil
}

func (r *Repository) GetArticleBySlug(slug string) (*Article, error) {
	query := `SELECT id, slug, title, summary, content, author, category, cover_image, is_published, p2p_signed, created_at, updated_at FROM cms_articles WHERE slug = ?`
	row := r.db.QueryRow(query, slug)

	var a Article
	var pub, p2p int
	if err := row.Scan(&a.ID, &a.Slug, &a.Title, &a.Summary, &a.Content, &a.Author, &a.Category, &a.CoverImage, &pub, &p2p, &a.CreatedAt, &a.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	a.IsPublished = pub == 1
	a.P2PSigned = p2p == 1

	return &a, nil
}

func (r *Repository) UpsertArticle(a Article) error {
	now := timeutil.Now()
	a.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(a.Slug), " ", "-"))
	if a.Slug == "" {
		a.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(a.Title), " ", "-"))
	}

	pubInt := 0
	if a.IsPublished {
		pubInt = 1
	}
	p2pInt := 0
	if a.P2PSigned {
		p2pInt = 1
	}

	_, err := r.db.Exec(`
		INSERT INTO cms_articles (id, slug, title, summary, content, author, category, cover_image, is_published, p2p_signed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			slug = excluded.slug,
			title = excluded.title,
			summary = excluded.summary,
			content = excluded.content,
			author = excluded.author,
			category = excluded.category,
			cover_image = excluded.cover_image,
			is_published = excluded.is_published,
			p2p_signed = excluded.p2p_signed,
			updated_at = excluded.updated_at
	`, a.ID, a.Slug, a.Title, a.Summary, a.Content, a.Author, a.Category, a.CoverImage, pubInt, p2pInt, now, now)

	return err
}

func (r *Repository) DeleteArticle(id string) error {
	_, err := r.db.Exec(`DELETE FROM cms_articles WHERE id = ?`, id)
	return err
}
