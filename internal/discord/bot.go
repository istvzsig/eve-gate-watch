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
	// if msgContent == IntelCommand || strings.HasPrefix(msgContent, IntelCommand+" ") {
	// 	dispatchIntelCommand(session, message, intelService)
	// 	return

	// }

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

func dispatchDscanCommand(
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
	security := formatSecurity(result.SecurityStatus)

	return fmt.Sprintf(
		"```ansi\n"+
			"%s SYSTEM INTELLIGENCE\n"+
			"--------------------------------\n"+
			"Security           %s\n"+
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
		security,
		result.Planets,
		result.Stargates,
		result.Stations,
		result.ShipJumps,
		result.ShipKills,
		result.PodKills,
		result.NPCKills,
	)
}

func formatSecurity(value float64) string {
	display := fmt.Sprintf("%.1f", value)

	switch {
	case value >= 0.5:
		return fmt.Sprintf("\033[32m%s HIGH SEC\033[0m", display)
	case value > 0.0:
		return fmt.Sprintf("\033[33m%s LOW SEC\033[0m", display)
	default:
		return fmt.Sprintf("\033[31m%s NULL SEC\033[0m", display)
	}
}

// func dispatchIntelCommand(
// 	session *discordgo.Session,
// 	message *discordgo.MessageCreate,
// 	intelService *intel.Service,
// ) {
// 	systemName := strings.TrimSpace(
// 		strings.TrimPrefix(message.Content, IntelCommand),
// 	)

// 	recentIntel, err := intelService.RecentIntel(
// 		context.Background(),
// 		systemName,
// 	)
// 	if err != nil {
// 		log.Printf("intel query: %v", err)
// 		return
// 	}

// 	result := formatIntel(recentIntel)

// 	_, err = session.ChannelMessageSend(
// 		message.ChannelID,
// 		result,
// 	)
// 	if err != nil {
// 		log.Printf("send intel message: %v", err)
// 	}
// }
