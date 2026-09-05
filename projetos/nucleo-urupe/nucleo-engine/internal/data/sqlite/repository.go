/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"nucleo-engine/internal/pkg/timeutil"
)

// Repository handles all SQLite operations for the V6 runtime.
// Logic is decomposed across several files in the same package:
// - message_repo.go: Message logging and reactions
// - memory_repo.go: Episodic memory and semantic search
// - budget_repo.go: Social budget and reservations
// - episode_repo.go: Episode management
// - identity_repo.go: Stance and behavior events
// - metrics_repo.go: Dashboard and system metrics
// - sync_repo.go: History synchronization state
// - persona_v33.go: Core identity and persona profiles
// - turn_repo.go: Turn requests and snapshots
type Repository struct {
	db *sql.DB
}

type executor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func (r *Repository) getExecutor(tx *sql.Tx) executor {
	if tx != nil {
		return tx
	}
	return r.db
}

// NewRepository initializes a new SQLite connection and applies the schema.
func NewRepository(dbPath string) (*Repository, error) {
	// If path is a simple filename, move it to the data/ directory
	if !strings.Contains(dbPath, "/") && !strings.Contains(dbPath, "\\") {
		dbPath = filepath.Join("data", dbPath)
	}

	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database at %s: %w", dbPath, err)
	}

	// Optimize for concurrency
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Apply performance pragmas and check errors
	pragmas := []string{
		"PRAGMA journal_mode=WAL",
		"PRAGMA synchronous=NORMAL",
		"PRAGMA busy_timeout=5000",
		"PRAGMA foreign_keys=ON",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, fmt.Errorf("failed to apply pragma %q: %w", p, err)
		}
	}

	repo := &Repository{db: db}
	if err := repo.initSchema(); err != nil {
		return nil, err
	}

	return repo, nil
}

// GetDB returns the underlying sql.DB.
func (r *Repository) GetDB() *sql.DB {
	return r.db
}

// Close closes the underlying sql.DB connection.
func (r *Repository) Close() error {
	return r.db.Close()
}

// GetConfig retrieves a value from system_config.
func (r *Repository) GetConfig(key string, defaultValue string) string {
	var value string
	err := r.db.QueryRow("SELECT value FROM system_config WHERE key = ?", key).Scan(&value)
	if err != nil {
		return defaultValue
	}
	return value
}

// SetConfig updates a value in system_config.
func (r *Repository) SetConfig(key string, value string) error {
	now := timeutil.Now()
	_, err := r.db.Exec(`
		INSERT INTO system_config (key, value, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(key) DO UPDATE SET
			value = excluded.value,
			updated_at = excluded.updated_at
	`, key, value, now)
	return err
}

// BeginImmediate starts a transaction with IMMEDIATE locking.
// In SQLite, this prevents deadlocks by acquiring a reserved lock immediately.
func (r *Repository) BeginImmediate() (*sql.Tx, error) {
	return r.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
}

// initSchema is implemented in migrations.go

// Helpers

// decodeStringSlice unmarshals a JSON array string into a slice of strings.
// It returns nil if the input is empty or invalid.
func decodeStringSlice(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

// uniqueAuthorIDs returns a slice of unique author IDs from the input slice.
// If limit > 0, it stops after reaching the specified count.
func uniqueAuthorIDs(authorIDs []string, limit int) []string {
	seen := make(map[string]bool, len(authorIDs))
	unique := make([]string, 0, len(authorIDs))
	for _, authorID := range authorIDs {
		authorID = strings.TrimSpace(authorID)
		if authorID == "" || seen[authorID] {
			continue
		}
		seen[authorID] = true
		unique = append(unique, authorID)
		if limit > 0 && len(unique) >= limit {
			break
		}
	}
	return unique
}
