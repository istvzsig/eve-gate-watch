package discord

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/iszsig/eve-gate-watch/internal/intel"
)

const (
	PingCommand  = "!ping"
	IntelCommand = "!intel"
	DScanCommand = "!dscan"
)

func CreateSession(token string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func OpenSession(session *discordgo.Session) error {
	if err := session.Open(); err != nil {
		return err
	}

	return nil
}

func CloseSession(session *discordgo.Session) error {
	if err := session.Close(); err != nil {
		return err
	}

	return nil
}

func RegisterHandlers(
	session *discordgo.Session,
	intelService *intel.Service,
) {
	session.AddHandler(func(
		session *discordgo.Session,
		message *discordgo.MessageCreate,
	) {
		dispatchCommands(session, message, intelService)
	})
}

func dispatchCommands(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
	intelService *intel.Service,
) {
	if message.Author.Bot {
		return
	}

	msgContent := strings.TrimSpace(message.Content)

	// ping
	if msgContent == PingCommand {
		dispatchPingCommand(session, message)
		return
	}

	// intel
	if msgContent == IntelCommand || strings.HasPrefix(msgContent, IntelCommand+" ") {
		dispatchIntelCommand(session, message, intelService)
		return

	}

	// dscan
	if msgContent == DScanCommand || strings.HasPrefix(msgContent, DScanCommand+" ") {
		dispatchDscanCommand(session, message, intelService)
		return
	}

}

func dispatchPingCommand(session *discordgo.Session, message *discordgo.MessageCreate) {
	_, err := session.ChannelMessageSend(
		message.ChannelID,
		"Pong! Gatewatch online.",
	)
	if err != nil {
		log.Printf("send message: %v", err)
	}
}

func dispatchIntelCommand(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
	intelService *intel.Service,
) {
	systemName := strings.TrimSpace(
		strings.TrimPrefix(message.Content, IntelCommand),
	)

	recentIntel, err := intelService.RecentIntel(
		context.Background(),
		systemName,
	)
	if err != nil {
		log.Printf("intel query: %v", err)
		return
	}

	result := formatIntel(recentIntel)

	_, err = session.ChannelMessageSend(
		message.ChannelID,
		result,
	)
	if err != nil {
		log.Printf("send intel message: %v", err)
	}
}

func dispatchDscanCommand(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
	intelService *intel.Service,
) {
	systemName := strings.TrimSpace(
		strings.TrimPrefix(message.Content, DScanCommand),
	)

	result, err := intelService.DScan(
		context.Background(),
		systemName,
	)
	if err != nil {
		log.Printf("dscan: %v", err)
		return
	}

	var builder strings.Builder

	fmt.Fprintf(
		&builder,
		"🚨 %s DSCAN\n\n",
		strings.ToUpper(result.SystemName),
	)

	for _, ship := range result.Ships {
		fmt.Fprintf(
			&builder,
			"🚀 %s ship\n",
			strings.ToUpper(ship.Name),
		)
	}

	_, err = session.ChannelMessageSend(
		message.ChannelID,
		builder.String(),
	)
	if err != nil {
		log.Printf("send dscan result: %v", err)
	}
}

func formatIntel(report intel.IntelReport) string {
	var builder strings.Builder

	fmt.Fprintf(
		&builder,
		"🛰️ %s INTEL\n\n",
		strings.ToUpper(report.System),
	)

	if len(report.Observations) == 0 {
		builder.WriteString("No recent intel.")
		return builder.String()
	}

	for _, observation := range report.Observations {
		fmt.Fprintf(
			&builder,
			"🚀 %s — %s\n",
			strings.ToUpper(observation.ShipName),
			observation.Threat,
		)
	}

	return builder.String()
}
