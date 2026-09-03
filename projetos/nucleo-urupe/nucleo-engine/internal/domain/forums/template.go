/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */

package forums

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"nucleo-engine/internal/data/sqlite"
)

type RenderedPost struct {
	Title string
	Body  string
	Tags  []string
}

func RenderTemplate(t sqlite.ForumTemplate, vars map[string]string) RenderedPost {
	title := substitute(t.Title, t.Variables, vars)
	body := substitute(t.Body, t.Variables, vars)
	return RenderedPost{Title: title, Body: body, Tags: t.Tags}
}

func substitute(text string, declared []string, vars map[string]string) string {
	result := text

	declaredMap := make(map[string]bool, len(declared))
	for _, v := range declared {
		declaredMap[v] = true
	}

	for _, name := range declared {
		val, ok := vars[name]
		if !ok {
			val = fmt.Sprintf("{{%s}}", name)
		}
		result = strings.ReplaceAll(result, "{{"+name+"}}", val)
	}

	if len(vars) == 0 {
		for _, name := range declared {
			val := builtinValue(name)
			result = strings.ReplaceAll(result, "{{"+name+"}}", val)
		}
	}

	return result
}

func builtinValue(name string) string {
	now := time.Now()
	switch strings.ToLower(name) {
	case "date":
		return now.Format("2006-01-02")
	case "time":
		return now.Format("15:04")
	case "datetime":
		return now.Format("2006-01-02 15:04")
	case "weekday":
		return now.Weekday().String()
	case "month":
		return now.Month().String()
	case "year":
		return fmt.Sprintf("%d", now.Year())
	case "random":
		return fmt.Sprintf("%d", rand.Intn(10000))
	default:
		return fmt.Sprintf("{{%s}}", name)
	}
}

func ShouldRunNow(t sqlite.ForumTemplate, lastRun time.Time) bool {
	now := time.Now()
	loc := now.Location()

	switch t.Schedule {
	case "daily":
		h, m := parseTime(t.ScheduleConfig)
		target := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, loc)
		return now.After(target) && lastRun.Before(target)

	case "weekly":
		h, m := parseTime(t.ScheduleConfig)
		day := parseWeekday(t.ScheduleConfig)
		daysUntil := (int(day) - int(now.Weekday()) + 7) % 7
		target := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, loc).AddDate(0, 0, daysUntil)
		return now.After(target) && lastRun.Before(target)

	case "monthly":
		h, m := parseTime(t.ScheduleConfig)
		dom := parseDayOfMonth(t.ScheduleConfig)
		target := time.Date(now.Year(), now.Month(), dom, h, m, 0, 0, loc)
		return now.After(target) && lastRun.Before(target)

	default:
		return false
	}
}

func parseTime(cfg map[string]any) (int, int) {
	h := 9
	m := 0
	if v, ok := cfg["hour"].(float64); ok {
		h = int(v)
	}
	if v, ok := cfg["minute"].(float64); ok {
		m = int(v)
	}
	return h, m
}

func parseWeekday(cfg map[string]any) time.Weekday {
	if v, ok := cfg["day"].(string); ok {
		switch strings.ToLower(v) {
		case "sunday":
			return time.Sunday
		case "monday":
			return time.Monday
		case "tuesday":
			return time.Tuesday
		case "wednesday":
			return time.Wednesday
		case "thursday":
			return time.Thursday
		case "friday":
			return time.Friday
		case "saturday":
			return time.Saturday
		}
	}
	return time.Monday
}

func parseDayOfMonth(cfg map[string]any) int {
	if v, ok := cfg["day"].(float64); ok {
		return int(v)
	}
	return 1
}
