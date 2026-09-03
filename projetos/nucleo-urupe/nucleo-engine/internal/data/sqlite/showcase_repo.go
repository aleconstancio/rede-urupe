package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type ShowcaseProject struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	Slug            string   `json:"slug"`
	Description     string   `json:"description"`
	LongDescription string   `json:"long_description"`
	Icon            string   `json:"icon"`
	URL             string   `json:"url"`
	GithubURL       string   `json:"github_url"`
	Category        string   `json:"category"`
	Status          string   `json:"status"`
	SortOrder       int      `json:"sort_order"`
	Tags            []string `json:"tags"`
	IsFeatured      bool     `json:"is_featured"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

func (r *Repository) ListShowcaseProjects() ([]ShowcaseProject, error) {
	rows, err := r.db.Query(`
		SELECT id, name, slug, description, long_description, icon, url, github_url,
		       category, status, sort_order, tags_json, is_featured, created_at, updated_at
		FROM showcase_projects
		ORDER BY sort_order ASC, name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list showcase projects: %w", err)
	}
	defer rows.Close()

	var projects []ShowcaseProject
	for rows.Next() {
		var p ShowcaseProject
		var tagsJSON string
		err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.LongDescription,
			&p.Icon, &p.URL, &p.GithubURL, &p.Category, &p.Status, &p.SortOrder,
			&tagsJSON, &p.IsFeatured, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan showcase project: %w", err)
		}
		json.Unmarshal([]byte(tagsJSON), &p.Tags)
		projects = append(projects, p)
	}
	return projects, nil
}

func (r *Repository) GetShowcaseProject(slug string) (*ShowcaseProject, error) {
	var p ShowcaseProject
	var tagsJSON string
	err := r.db.QueryRow(`
		SELECT id, name, slug, description, long_description, icon, url, github_url,
		       category, status, sort_order, tags_json, is_featured, created_at, updated_at
		FROM showcase_projects WHERE slug = ?
	`, slug).Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.LongDescription,
		&p.Icon, &p.URL, &p.GithubURL, &p.Category, &p.Status, &p.SortOrder,
		&tagsJSON, &p.IsFeatured, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get showcase project: %w", err)
	}
	json.Unmarshal([]byte(tagsJSON), &p.Tags)
	return &p, nil
}

func (r *Repository) CreateShowcaseProject(p *ShowcaseProject) error {
	tagsJSON, _ := json.Marshal(p.Tags)
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	result, err := r.db.Exec(`
		INSERT INTO showcase_projects (name, slug, description, long_description, icon, url, github_url,
		                              category, status, sort_order, tags_json, is_featured, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, p.Name, p.Slug, p.Description, p.LongDescription, p.Icon, p.URL, p.GithubURL,
		p.Category, p.Status, p.SortOrder, string(tagsJSON), p.IsFeatured, now, now)
	if err != nil {
		return fmt.Errorf("create showcase project: %w", err)
	}
	p.ID, _ = result.LastInsertId()
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

func (r *Repository) UpdateShowcaseProject(p *ShowcaseProject) error {
	tagsJSON, _ := json.Marshal(p.Tags)
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	_, err := r.db.Exec(`
		UPDATE showcase_projects SET name=?, description=?, long_description=?, icon=?, url=?,
		       github_url=?, category=?, status=?, sort_order=?, tags_json=?, is_featured=?, updated_at=?
		WHERE slug=?
	`, p.Name, p.Description, p.LongDescription, p.Icon, p.URL, p.GithubURL,
		p.Category, p.Status, p.SortOrder, string(tagsJSON), p.IsFeatured, now, p.Slug)
	if err != nil {
		return fmt.Errorf("update showcase project: %w", err)
	}
	return nil
}

func (r *Repository) DeleteShowcaseProject(slug string) error {
	_, err := r.db.Exec("DELETE FROM showcase_projects WHERE slug = ?", slug)
	if err != nil {
		return fmt.Errorf("delete showcase project: %w", err)
	}
	return nil
}
