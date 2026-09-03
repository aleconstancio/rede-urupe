/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package postgres

import (
	"database/sql"
	"time"

	"github.com/aleconstancio/minos/db"
)

// Connect opens a Postgres connection pool using Talos's db.Connect.
// The DSN should be a Postgres connection string:
//
//	postgres://user:pass@host:5432/dbname?sslmode=require
//
// This package is not wired into Maze's startup yet — it exists alongside
// the SQLite repository for the Phase B migration. When ready, swap the
// import in cmd/bot/main.go and adjust queries for Postgres syntax.
func Connect(dsn string) (*sql.DB, error) {
	return db.Connect(dsn,
		db.WithDriver("postgres"),
		db.WithPool(25, 5),
		db.WithConnLifetime(5*time.Minute, 3*time.Minute),
		db.WithRetry(10, 3*time.Second),
	)
}
