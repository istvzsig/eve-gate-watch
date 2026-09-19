package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is required")
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("create Discord session: %v", err)
	}

	session.AddHandler(handleMessage)

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
) {
	if message.Author.Bot {
		return
	}

	if message.Content != "!ping" {
		return
	}

	_, err := session.ChannelMessageSend(
		message.ChannelID,
		"Pong! Gatewatch online.",
	)
	if err != nil {
		log.Printf("send message: %v", err)
	}
}
