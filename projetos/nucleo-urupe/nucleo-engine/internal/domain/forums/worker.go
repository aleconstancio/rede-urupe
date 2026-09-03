/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */

// Package forums provides scheduled forum post publishing from templates.
package forums

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/pkg/timeutil"
)

type ForumWorker struct {
	repo *sqlite.Repository
	sess *discordgo.Session
}

func NewForumWorker(repo *sqlite.Repository) *ForumWorker {
	return &ForumWorker{repo: repo}
}

func (w *ForumWorker) SetSession(s *discordgo.Session) {
	w.sess = s
}

func (w *ForumWorker) Start(ctx context.Context) {
	log.Println("[ForumWorker] Started")
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	w.runScheduled(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("[ForumWorker] Shutting down")
			return
		case <-ticker.C:
			w.runScheduled(ctx)
		}
	}
}

func (w *ForumWorker) runScheduled(ctx context.Context) {
	templates, err := w.repo.GetScheduledForumTemplates()
	if err != nil {
		log.Printf("[ForumWorker] Failed to load scheduled templates: %v", err)
		return
	}

	for _, t := range templates {
		now := timeutil.Now()
		if !ShouldRunNow(t, now.Add(-24*time.Hour)) {
			continue
		}

		post, err := w.PublishFromTemplate(ctx, t, nil)
		if err != nil {
			log.Printf("[ForumWorker] Failed to publish scheduled post %q: %v", t.Title, err)
			continue
		}
		log.Printf("[ForumWorker] Published scheduled forum post %q (msg=%s)", post.Title, post.DiscordMessageID)
	}
}

func (w *ForumWorker) PublishFromTemplate(ctx context.Context, t sqlite.ForumTemplate, vars map[string]string) (*sqlite.ForumPost, error) {
	if w.sess == nil {
		return nil, fmt.Errorf("discord session not configured")
	}

	rendered := RenderTemplate(t, vars)

	post := sqlite.ForumPost{
		TemplateID: t.ID,
		GuildID:    t.GuildID,
		ChannelID:  t.ChannelID,
		Title:      rendered.Title,
		Body:       rendered.Body,
		Tags:       rendered.Tags,
		Status:     "draft",
	}

	// Check if channel is a forum channel
	ch, err := w.sess.Channel(t.ChannelID)
	if err != nil {
		post.Status = "failed"
		post.Error = err.Error()
		w.repo.SaveForumPost(post)
		return nil, fmt.Errorf("get channel: %w", err)
	}

	var availableTags []discordgo.ForumTag
	if ch.Type == discordgo.ChannelTypeGuildForum {
		availableTags = ch.AvailableTags
	}

	var tagIDs []string
	for _, tag := range rendered.Tags {
		for _, at := range availableTags {
			if strings.EqualFold(tag, at.Name) {
				tagIDs = append(tagIDs, at.ID)
				break
			}
		}
	}

	th, err := w.sess.ForumThreadStartComplex(t.ChannelID, &discordgo.ThreadStart{
		Name:        rendered.Title,
		AppliedTags: tagIDs,
	}, &discordgo.MessageSend{
		Content: rendered.Body,
	})
	if err != nil {
		post.Status = "failed"
		post.Error = err.Error()
		w.repo.SaveForumPost(post)
		return nil, fmt.Errorf("create forum post: %w", err)
	}

	threadID := th.ID

	post.DiscordThreadID = threadID
	post.Status = "published"

	id, err := w.repo.SaveForumPost(post)
	if err != nil {
		log.Printf("[ForumWorker] Failed to save forum post record: %v", err)
	}
	post.ID = id

	return &post, nil
}
