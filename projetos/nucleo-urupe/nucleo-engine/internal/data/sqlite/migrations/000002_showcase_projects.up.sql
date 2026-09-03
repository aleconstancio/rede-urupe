CREATE TABLE IF NOT EXISTS showcase_projects (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    long_description TEXT DEFAULT '',
    icon TEXT DEFAULT '',
    url TEXT DEFAULT '',
    github_url TEXT DEFAULT '',
    category TEXT DEFAULT 'product',
    status TEXT DEFAULT 'active',
    sort_order INTEGER DEFAULT 0,
    tags_json TEXT DEFAULT '[]',
    is_featured INTEGER DEFAULT 0,
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX idx_showcase_projects_slug ON showcase_projects(slug);
CREATE INDEX idx_showcase_projects_category ON showcase_projects(category);
CREATE INDEX idx_showcase_projects_sort ON showcase_projects(sort_order);
