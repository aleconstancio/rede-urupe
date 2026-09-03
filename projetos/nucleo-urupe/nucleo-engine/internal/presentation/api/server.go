/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"nucleo-engine/internal/config"
	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/domain/forums"
	"nucleo-engine/internal/pkg/timeutil"
)

type Server struct {
	repo        *sqlite.Repository
	cfg         config.Config
	forumWorker *forums.ForumWorker
	startedAt   time.Time

	// SSE Clients
	mu      sync.Mutex
	clients map[chan string]bool

	rateLimiter *RateLimiter
}

func NewServer(repo *sqlite.Repository, cfg config.Config) *Server {
	return &Server{
		repo:        repo,
		cfg:         cfg,
		startedAt:   time.Now(),
		clients:     make(map[chan string]bool),
		rateLimiter: NewRateLimiter(60, time.Minute),
	}
}

func (s *Server) SetForumWorker(w *forums.ForumWorker) {
	s.forumWorker = w
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/metrics", s.handleMetrics)
	mux.HandleFunc("/api/feed", s.handleFeed)
	mux.HandleFunc("/api/memory/today", s.handleMemoryToday)
	mux.HandleFunc("/api/topics", s.handleTopics) // Deprecated
	mux.HandleFunc("/api/persona", s.handlePersona)
	mux.HandleFunc("/api/persona/save", s.handlePersonaSave)
	mux.HandleFunc("/api/identities", s.handleIdentities)
	mux.HandleFunc("/api/identities/", s.handleIdentityByID)
	mux.HandleFunc("/api/personas", s.handlePersonas)
	mux.HandleFunc("/api/personas/", s.handlePersonaByID)
	mux.HandleFunc("/api/persona-policy", s.handlePersonaPolicy)
	mux.HandleFunc("/api/persona-adaptations", s.handlePersonaAdaptations)
	mux.HandleFunc("/api/persona-adaptations/reset", s.handlePersonaAdaptationsReset)
	mux.HandleFunc("/api/persona-proposals", s.handlePersonaProposals)
	mux.HandleFunc("/api/persona-proposals/", s.handlePersonaProposalAction)
	mux.HandleFunc("/api/mode", s.handleMode)
	mux.HandleFunc("/api/stats", s.handleStats)
	mux.HandleFunc("/api/annotations", s.handleAnnotations)
	mux.HandleFunc("/api/forum/templates", s.handleForumTemplates)
	mux.HandleFunc("/api/forum/publish", s.handleForumPublish)
	mux.HandleFunc("/api/forum/posts", s.handleForumPosts)
	mux.HandleFunc("/api/admin/members", s.requireAuth(s.handleAdminMembers))
	mux.HandleFunc("/api/admin/member/", s.requireAuth(s.handleAdminMemberByID))
	mux.HandleFunc("/api/admin/audit", s.requireAuth(s.handleAdminAudit))
	mux.HandleFunc("/api/admin/modlog", s.requireAuth(s.handleAdminModLog))
	mux.HandleFunc("/api/admin/welcome", s.requireAuth(s.handleAdminWelcome))
	mux.HandleFunc("/api/analytics/sentiment", s.handleAnalyticsSentiment)
	mux.HandleFunc("/api/analytics/tokens", s.handleAnalyticsTokens)
	mux.HandleFunc("/api/analytics/growth", s.handleAnalyticsGrowth)
	mux.HandleFunc("/api/analytics/channels", s.handleAnalyticsChannels)
	mux.HandleFunc("/api/analytics/engagement", s.handleAnalyticsEngagement)
	mux.HandleFunc("/events", s.rateLimiter.Middleware(s.handleEvents))
	mux.HandleFunc("/api/projects", s.handleShowcaseProjects)
	mux.HandleFunc("/api/projects/", s.handleShowcaseProjectBySlug)
	mux.HandleFunc("/api/channels", s.handleChannels)
	mux.HandleFunc("/api/cms/articles", handleCMSArticlesList(s.repo))
	mux.HandleFunc("/api/cms/article", handleCMSArticleGet(s.repo))
	mux.HandleFunc("/api/cms/article/save", handleCMSArticleSave(s.repo))
	mux.HandleFunc("/api/cms/article/delete", handleCMSArticleDelete(s.repo))

	// Serve Quartz knowledge map at /knowledge/*
	knowledgePath := "internal/presentation/web/knowledge_dist"
	if _, err := os.Stat(knowledgePath); err == nil {
		knowledgeFS := http.FileServer(http.Dir(knowledgePath))
		mux.Handle("/knowledge/", http.StripPrefix("/knowledge/", knowledgeFS))
		log.Printf("[API] Serving knowledge map from %s at /knowledge/", knowledgePath)
	} else {
		log.Printf("[API] Knowledge dist not found at %s. Skipping knowledge map.", knowledgePath)
	}

	// Serve frontend static files from "internal/presentation/web/static_dist"
	staticPath := "internal/presentation/web/static_dist"
	if _, err := os.Stat(staticPath); err == nil {
		mux.Handle("/", http.FileServer(http.Dir(staticPath)))
		log.Printf("[API] Serving static files from %s", staticPath)
	} else {
		log.Printf("[API] Static files not found at %s. Skipping dashboard UI.", staticPath)
	}

	addr := ":" + s.cfg.DashboardPort
	if addr == ":" {
		addr = ":8080"
	}

	log.Printf("[API] Dashboard server starting on %s", addr)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Printf("[CRITICAL] API Server failed: %v", err)
		}
	}()
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	uptimeMinutes := int(time.Since(s.startedAt).Minutes())

	totalMessages := 0
	s.repo.GetDB().QueryRow("SELECT COUNT(*) FROM messages").Scan(&totalMessages)

	activeMembers := 0
	weekAgo := timeutil.Now().Add(-7 * 24 * time.Hour)
	s.repo.GetDB().QueryRow("SELECT COUNT(DISTINCT author_id) FROM messages WHERE timestamp > ?", weekAgo).Scan(&activeMembers)

	botMessagesToday := 0
	today := timeutil.Today()
	s.repo.GetDB().QueryRow("SELECT COUNT(*) FROM messages WHERE is_bot = 1 AND timestamp > ?", today).Scan(&botMessagesToday)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"online":             true,
		"uptime_minutes":     uptimeMinutes,
		"total_messages":     totalMessages,
		"active_members_7d":  activeMembers,
		"bot_messages_today": botMessagesToday,
		"channel_name":       s.cfg.TargetChannelName,
		"version":            "2.0.0-airelius",
	})
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	channelID := s.cfg.TargetChannelID

	// 1. Total Cost (24h)
	totalCost, _ := s.repo.GetTotalCostLast24Hours()

	// 2. Active Users (1h)
	activeUsers := 0
	hourAgo := timeutil.Now().Add(-time.Hour)
	s.repo.GetDB().QueryRow("SELECT COUNT(DISTINCT author_id) FROM messages WHERE channel_id = ? AND timestamp > ?", channelID, hourAgo).Scan(&activeUsers)

	// 3. Messages (1h)
	messages1h := 0
	s.repo.GetDB().QueryRow("SELECT COUNT(*) FROM messages WHERE channel_id = ? AND timestamp > ?", channelID, hourAgo).Scan(&messages1h)

	// 4. Bot Messages (Today)
	botMsgs := 0
	today := timeutil.Today()
	s.repo.GetDB().QueryRow("SELECT COUNT(*) FROM messages WHERE channel_id = ? AND is_bot = 1 AND timestamp > ?", channelID, today).Scan(&botMsgs)

	// 5. Memories (Today)
	memories := 0
	s.repo.GetDB().QueryRow("SELECT COUNT(*) FROM memory_capsules WHERE created_at > ?", today).Scan(&memories)

	// 6. Active Persona
	personaName := "Bot"
	studio, err := s.repo.GetPersonaStudioState(channelID)
	if err == nil {
		personaName = studio.ActiveIdentity.Name
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"active_users":     activeUsers,
		"messages_1h":      messages1h,
		"channel_name":     s.cfg.TargetChannelName,
		"active_persona":   personaName,
		"bot_messages":     botMsgs,
		"memories_today":   memories,
		"total_tokens_24h": totalCost,
	})
}

func (s *Server) handleFeed(w http.ResponseWriter, r *http.Request) {
	msgs, _ := s.repo.GetRecentMessages(s.cfg.TargetChannelID, 20)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"Messages":   msgs,
		"CountLabel": fmt.Sprintf("%d msgs", len(msgs)),
		"EmptyText":  "O canal está em silêncio. Comece uma conversa para ver o pulso de Maze.",
	})
}

func (s *Server) handleMemoryToday(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = timeutil.Today()
	}

	capsules, err := s.repo.GetActiveMemoryCapsulesByDate(dateStr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"day":      dateStr,
		"capsules": capsules,
	})
}

func (s *Server) handleTopics(w http.ResponseWriter, r *http.Request) {
	// Deprecated in V5
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "deprecated",
		"note":   "Use /api/memory/today instead",
	})
}

func (s *Server) handlePersona(w http.ResponseWriter, r *http.Request) {
	channelID := s.personaChannelID(r)
	studioState, err := s.repo.GetPersonaStudioState(channelID)
	if err == nil {
		passiveMode := s.repo.GetConfig("passive_mode", "true") == "true"

		identities, _ := s.repo.ListCoreIdentityProfiles()
		overlays, _ := s.repo.ListPersonaOverlays("")

		json.NewEncoder(w).Encode(map[string]interface{}{
			"ActiveProfile": map[string]interface{}{
				"Name":        studioState.ActiveIdentity.Name,
				"PassiveMode": passiveMode,
			},
			"ActiveIdentity":   studioState.ActiveIdentity,
			"ActivePersona":    studioState.ActiveOverlay,
			"PersonaPolicy":    studioState.Policy,
			"AdaptiveMemories": studioState.AdaptiveMemories,
			"Identities":       identities,
			"Overlays":         overlays,
			"Capabilities": map[string]bool{
				"v5_gateway":     true,
				"multi_identity": true,
			},
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"ActiveProfile": map[string]interface{}{
			"Name": "Bot",
		},
	})
}

func (s *Server) handleMode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Passive bool `json:"passive"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	val := "false"
	if req.Passive {
		val = "true"
	}
	s.repo.SetConfig("passive_mode", val)
	s.Broadcast("persona")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"passive": req.Passive,
	})
}

func (s *Server) handlePersonaSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}
	s.handlePersona(w, r)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	// Simplified stats for V5
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "V6 AIrelius Agent Pipeline Active",
	})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	client := make(chan string)
	s.mu.Lock()
	s.clients[client] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, client)
		s.mu.Unlock()
		close(client)
	}()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-client:
			fmt.Fprintf(w, "event: %s\ndata: {}\n\n", msg)
			w.(http.Flusher).Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (s *Server) handleAnnotations(w http.ResponseWriter, r *http.Request) {
	channelID := s.cfg.TargetChannelID
	items, err := s.repo.ListMessageAnnotations(channelID, 50)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{"items": items})
}

func (s *Server) handleChannels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	channels, err := s.repo.ListActiveChannels(s.cfg.TargetGuildID)
	if err != nil {
		channels = []string{}
	}
	hasDefault := false
	for _, ch := range channels {
		if ch == s.cfg.TargetChannelID {
			hasDefault = true
			break
		}
	}
	if !hasDefault && s.cfg.TargetChannelID != "" {
		channels = append([]string{s.cfg.TargetChannelID}, channels...)
	}
	json.NewEncoder(w).Encode(map[string]any{"channels": channels})
}

func (s *Server) Broadcast(event string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for client := range s.clients {
		select {
		case client <- event:
		default:
		}
	}
}
