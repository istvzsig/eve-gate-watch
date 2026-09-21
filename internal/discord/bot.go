package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/iszsig/eve-gate-watch/internal/dscan"
	"github.com/iszsig/eve-gate-watch/internal/intel"
)

const (
	PingCommand  = "!ping"
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

	// dscan
	if msgContent == DScanCommand || strings.HasPrefix(msgContent, DScanCommand+" ") {
		dispatchDScanCommand(session, message, intelService)
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

func dispatchDScanCommand(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
	intelService *intel.Service,
) {
	systemName := strings.TrimSpace(
		strings.TrimPrefix(message.Content, DScanCommand),
	)

	result, err := intelService.DScan(systemName)
	if err != nil {
		log.Printf("dscan: %v", err)

		_, sendErr := session.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprintf(
				"System `%s` not found.",
				systemName,
			),
		)
		if sendErr != nil {
			log.Printf("send dscan error: %v", sendErr)
		}

		return
	}

	messageText := formatDScan(result)

	_, err = session.ChannelMessageSend(
		message.ChannelID,
		messageText,
	)
	if err != nil {
		log.Printf("send dscan result: %v", err)
	}
}

func formatDScan(result dscan.DScanResult) string {
	return fmt.Sprintf(
		"```text\n"+
			"%s SYSTEM INTELLIGENCE\n"+
			"--------------------------------\n"+
			"Region             %s\n"+
			"Constellation      %s\n"+
			"Security           %.1f %s\n"+
			"Planets            %d\n"+
			"Stargates          %d\n"+
			"Stations           %d\n"+
			"\n"+
			"LAST HOUR\n"+
			"--------------------------------\n"+
			"Ship jumps         %d\n"+
			"Ship kills         %d\n"+
			"Pod kills          %d\n"+
			"NPC kills          %d\n"+
			"```",
		strings.ToUpper(result.SystemName),
		result.RegionName,
		result.ConstellationName,
		result.SecurityStatus,
		securityClass(result.SecurityStatus),
		result.Planets,
		result.Stargates,
		result.Stations,
		result.ShipJumps,
		result.ShipKills,
		result.PodKills,
		result.NPCKills,
	)
}

func securityClass(value float64) string {
	switch {
	case value >= 0.5:
		return "HIGH SEC"
	case value > 0.0:
		return "LOW SEC"
	default:
		return "NULL SEC"
	}
}
