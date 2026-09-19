package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/iszsig/eve-gate-watch/internal/intel"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is required")
	}

	repository := intel.NewMemoryRepository()
	intelService := intel.NewService(repository)

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("create Discord session: %v", err)
	}

	session.AddHandler(func(
		session *discordgo.Session,
		message *discordgo.MessageCreate,
	) {
		handleMessage(session, message, intelService)
	})

	if err := session.Open(); err != nil {
		log.Fatalf("open Discord session: %v", err)
	}

	log.Println("Gatewatch is online")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("shutting down")

	if err := session.Close(); err != nil {
		log.Printf("close Discord session: %v", err)
	}
}

func handleMessage(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
	intelService *intel.Service,
) {
	if message.Author.Bot {
		return
	}

	content := strings.TrimSpace(message.Content)

	if content == "!ping" {
		_, err := session.ChannelMessageSend(
			message.ChannelID,
			"Pong! Gatewatch online.",
		)
		if err != nil {
			log.Printf("send message: %v", err)
		}

		return
	}

	if strings.HasPrefix(content, "!intel") {
		systemName := strings.TrimSpace(
			strings.TrimPrefix(content, "!intel"),
		)

		result, err := intelService.RecentIntel(
			context.Background(),
			systemName,
		)
		if err != nil {
			log.Printf("intel query: %v", err)
			return
		}

		_, err = session.ChannelMessageSend(
			message.ChannelID,
			result,
		)
		if err != nil {
			log.Printf("send intel message: %v", err)
		}

		return
	}
}
