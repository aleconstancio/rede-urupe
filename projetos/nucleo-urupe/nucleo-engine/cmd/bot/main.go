/*
 * Copyright (c) 2026 Maze Project
 * Licensed under the MIT License. See LICENSE in the project root for license information.
 */
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aleconstancio/talos/v2/engine/core/agent"
	minotaur "github.com/aleconstancio/talos/v2/behavior/persona"
	"github.com/bwmarrin/discordgo"

	"nucleo-engine/internal/config"
	"nucleo-engine/internal/data/llm"
	"nucleo-engine/internal/data/sqlite"
	"nucleo-engine/internal/domain/commands"
	"nucleo-engine/internal/domain/forums"
	"nucleo-engine/internal/domain/gateway"
	"nucleo-engine/internal/domain/memory"
	"nucleo-engine/internal/presentation/api"
	"nucleo-engine/internal/presentation/discord"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := sqlite.NewRepository(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite: %v", err)
	}

	llmClient := llm.NewClient(cfg, repo, repo)

	// Phase 2: Gateway Initialization
	gtw := gateway.NewGateway(llmClient)
	gatewayWorker := gateway.NewGatewayWorker(repo, gtw, cfg.TargetGuildID, cfg.TargetChannelID, cfg.GatewayGateModel, cfg.GatewayReplyModel)

	// Phase 2b: Talos Agent Pipeline
	assembler := gateway.NewPayloadAssembler(repo)
	resolver := minotaur.NewResolver(repo)
	messenger := gateway.NewDiscordMessenger(nil, repo, nil)

	mazeIngester := &gateway.TriggerIngester{Repo: repo}
	mazeContextAssembler := &gateway.MazeContextAssembler{Repo: repo, Assembler: assembler}
	mazePersonaSelector := &gateway.MazePersonaSelector{Resolver: resolver}
	mazeGater := &gateway.MazeGater{Gateway: gtw, Assembler: assembler, GateModel: cfg.GatewayGateModel}
	mazeGenerator := &gateway.MazeGenerator{Gateway: gtw, Messenger: messenger, Assembler: assembler, Repo: repo, Cfg: cfg}
	pulseDetector := &gateway.MazePulseDetector{Repo: repo}
	backlogStore := &gateway.SQLiteBacklogStore{Repo: repo}

	agentPipeline := agent.NewAgent(
		agent.WithIngester(mazeIngester),
		agent.WithContextAssembler(mazeContextAssembler),
		agent.WithPersonaSelector(mazePersonaSelector),
		agent.WithGater(mazeGater),
		agent.WithGenerator(mazeGenerator),
	)

	agentWorker := gateway.NewAgentWorker(agentPipeline, repo, pulseDetector, backlogStore, nil)

	// Phase 3: Fractal Memory Workers
	capsuleWorker := memory.NewCapsuleWorker(repo, llmClient, cfg.MemoryModel, cfg.TargetChannelID)
	compactor := memory.NewCompactor(repo, llmClient, cfg.MemoryModel)

	// Phase 4: Forum Automation
	forumWorker := forums.NewForumWorker(repo)

	handler := discord.NewHandler(gatewayWorker.TriggerChan, repo, cfg.TargetGuildID, cfg.TargetChannelID)
	handler.SetGatewayWorker(gatewayWorker)
	handler.SetMultiChannel(cfg.MultiChannel)

	cmdRouter := commands.NewRouter()
	cmdRouter.Register(commands.StatusCommand(repo))
	cmdRouter.Register(commands.ConfigCommand(repo))
	cmdRouter.Register(commands.MemberCommand(repo))
	cmdRouter.Register(commands.AdminCommand(repo))
	handler.SetCommandRouter(cmdRouter)

	if cfg.DashboardEnabled {
		// Update API server initialization to remove orchestrator
		apiServer := api.NewServer(repo, cfg)
		apiServer.SetForumWorker(forumWorker)
		apiServer.Start()
		gatewayWorker.SetBroadcaster(apiServer)
		agentWorker.SetBroadcaster(apiServer)
	}

	_ = agentWorker // Available for future Agent-pipeline event dispatch

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start background workers
	go gatewayWorker.Start(ctx)
	go capsuleWorker.Start(ctx)
	go compactor.Start(ctx)
	go forumWorker.Start(ctx)

	if cfg.DiscordToken == "" {
		log.Printf("[AVISO] DISCORD_TOKEN não configurado. Rodando em Modo Local Offline (apenas servidor Web/Dashboard em http://localhost:%s).", cfg.DashboardPort)
		log.Printf("[INFO] Para conectar a Micélia 🍄 ao Discord, adicione DISCORD_TOKEN e TARGET_GUILD_ID no seu arquivo .env.")
		
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
		log.Println("[INFO] Encerrando servidor local do Núcleo Urupê.")
		return
	}

	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentMessageContent | discordgo.IntentsGuildMembers

	dg.AddHandler(handler.OnMessageCreate)
	dg.AddHandler(handler.OnMessageUpdate)
	dg.AddHandler(handler.OnMessageDelete)
	dg.AddHandler(handler.OnReactionAdd)
	dg.AddHandler(handler.OnReactionRemove)
	dg.AddHandler(handler.OnReactionRemoveAll)
	dg.AddHandler(handler.OnInteractionCreate)
	dg.AddHandler(handler.OnGuildMemberAdd)
	dg.AddHandler(handler.OnGuildMemberRemove)

	if err := dg.Open(); err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer dg.Close()

	// Resolve channel name for dashboard
	if ch, err := dg.Channel(cfg.TargetChannelID); err == nil {
		cfg.TargetChannelName = "#" + ch.Name
	} else {
		cfg.TargetChannelName = "Canal Desconhecido"
	}

	gatewayWorker.SetSession(dg)
	forumWorker.SetSession(dg)
	handler.Start(dg)

	log.Printf("Núcleo Urupê (Micélia 🍄 & Gateway) rodando no canal %s. Gate: %s. Reply: %s. Memória: %s.", cfg.TargetChannelID, cfg.GatewayGateModel, cfg.GatewayReplyModel, cfg.MemoryModel)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shutting down...")
	handler.Close()
}
