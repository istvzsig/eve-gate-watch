package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/istvzsig/eve-gate-watch/external/esi"
	"github.com/istvzsig/eve-gate-watch/external/zkill"
	"github.com/istvzsig/eve-gate-watch/internal/data"
	"github.com/istvzsig/eve-gate-watch/internal/discord"
	"github.com/istvzsig/eve-gate-watch/internal/dscan"
	"github.com/istvzsig/eve-gate-watch/internal/httpapi"
	"github.com/istvzsig/eve-gate-watch/internal/intel"
	"github.com/istvzsig/eve-gate-watch/internal/systems"

	_ "github.com/lib/pq"
)

const (
	esiBaseURL = "https://esi.evetech.net/latest"
	zkillURL   = "https://zkillboard.com/api"
)

var _ data.Repository = (*data.PostgresRepository)(nil)
var _ dscan.Scanner = (*esi.ESIClient)(nil)
var _ zkill.SystemResolver = (*esi.ESIClient)(nil)
var _ zkill.TypeResolver = (*esi.ESIClient)(nil)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is required")
	}

	esiClient := esi.NewESIClient(esiBaseURL)
	zkillClient := zkill.NewClient(zkillURL, esiClient, esiClient)

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("database connection: %v", err)
	}

	postgresRepository := data.NewPostgresRepository(db)

	intelService := intel.NewService(esiClient, zkillClient)
	systemsService := systems.NewService(postgresRepository)

	// Http api
	go func() {
		mux := http.NewServeMux()
		httpHandler := httpapi.NewHandler(intelService, systemsService)

		mux.HandleFunc("/api/dscan", httpHandler.HandleDScan)
		mux.HandleFunc("/api/pvp", httpHandler.HandleRecentPVP)
		mux.HandleFunc("/api/systems", httpHandler.HandleSystems)

		log.Println("HTTP server listening on :8080")

		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// // Discord bot session
	botSession, err := discord.CreateSession(token)
	if err != nil {
		log.Fatalf("create Discord session: %v", err)
		return
	}

	if err := discord.OpenSession(botSession); err != nil {
		log.Fatalf("open Discord session: %v", err)
		return
	}

	discord.RegisterPrefixCommands(botSession, intelService, systemsService)
	discord.RegisterSlashCommands(botSession)

	log.Println("Gatewatch bot is online")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("shutting down discord session")

	if err := discord.CloseSession(botSession); err != nil {
		log.Printf("close Discord session: %v", err)
	}
}
