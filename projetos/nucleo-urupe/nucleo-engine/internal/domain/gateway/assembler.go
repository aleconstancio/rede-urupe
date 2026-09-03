/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package gateway

import (
	"strings"

	"nucleo-engine/internal/data/sqlite"
)

type PayloadAssembler struct {
	repo   *sqlite.Repository
	emojis map[string]string // Name -> ID
}

func NewPayloadAssembler(repo *sqlite.Repository) *PayloadAssembler {
	return &PayloadAssembler{
		repo:   repo,
		emojis: make(map[string]string),
	}
}

func (a *PayloadAssembler) SetEmojis(emojis map[string]string) {
	a.emojis = emojis
}


// IsImageURL returns true if the URL points to a supported image format.
func IsImageURL(url string) bool {
	u := strings.ToLower(url)
	if idx := strings.Index(u, "?"); idx != -1 {
		u = u[:idx]
	}
	return strings.HasSuffix(u, ".jpg") ||
		strings.HasSuffix(u, ".jpeg") ||
		strings.HasSuffix(u, ".png") ||
		strings.HasSuffix(u, ".webp") ||
		strings.HasSuffix(u, ".gif")
}
