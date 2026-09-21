package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/istvzsig/eve-gate-watch/external/esi"
	"github.com/istvzsig/eve-gate-watch/external/zkill"
	"github.com/istvzsig/eve-gate-watch/internal/discord"
	"github.com/istvzsig/eve-gate-watch/internal/dscan"
	"github.com/istvzsig/eve-gate-watch/internal/httpapi"
	"github.com/istvzsig/eve-gate-watch/internal/intel"
)

const (
	esiBaseURL = "https://esi.evetech.net/latest"
	zkillURL   = "https://zkillboard.com/api"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is required")
	}

	esiClient := esi.NewESIClient(esiBaseURL)

	var _ dscan.Scanner = (*esi.ESIClient)(nil)
	var _ zkill.SystemResolver = (*esi.ESIClient)(nil)
	var _ zkill.TypeResolver = (*esi.ESIClient)(nil)

	zkillClient := zkill.NewClient(zkillURL, esiClient, esiClient)

	intelService := intel.NewService(esiClient, zkillClient)

	// Http api
	go func() {
		mux := http.NewServeMux()
		handler := httpapi.NewHandler(intelService)

		mux.HandleFunc("/api/dscan", handler.HandleDScan)
		mux.HandleFunc("/api/pvp", handler.HandleRecentPVP)

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

	discord.RegisterHandlers(botSession, intelService)

	if err := discord.OpenSession(botSession); err != nil {
		log.Fatalf("open Discord session: %v", err)
		return
	}

	log.Println("Gatewatch bot is online")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("shutting down discord session")

	if err := discord.CloseSession(botSession); err != nil {
		log.Printf("close Discord session: %v", err)
	}
}
