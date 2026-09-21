package main

import (
	"log"
	"net/http"
	"os"

	"github.com/iszsig/eve-gate-watch/internal/dscan"
	"github.com/iszsig/eve-gate-watch/internal/httpapi"
	"github.com/iszsig/eve-gate-watch/internal/intel"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is required")
	}

	repository := intel.NewMemoryRepository()

	var dScanner dscan.Scanner

	intelService := intel.NewService(repository, dScanner)

	// Http api
	go func() {
		mux := http.NewServeMux()
		handler := httpapi.NewHandler(intelService)

		mux.HandleFunc("/api/dscan", handler.HandleDScan)
		mux.HandleFunc("/api/intel", handler.HandleIntel)

		log.Println("HTTP server listening on :8080")

		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	for {
		select {}
	}

	// // Discord bot session
	// botSession, err := discord.CreateSession(token)
	// if err != nil {
	// 	log.Fatalf("create Discord session: %v", err)
	// 	return
	// }

	// discord.RegisterHandlers(botSession, intelService)

	// if err := discord.OpenSession(botSession); err != nil {
	// 	log.Fatalf("open Discord session: %v", err)
	// 	return
	// }

	// log.Println("Gatewatch bot is online")

	// stop := make(chan os.Signal, 1)
	// signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// <-stop

	// log.Println("shutting down discord session")

	// if err := discord.CloseSession(botSession); err != nil {
	// 	log.Printf("close Discord session: %v", err)

	// }

}
