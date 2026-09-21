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
	var builder strings.Builder

	fmt.Fprintf(
		&builder,
		"🛰️ %s INTEL\n\n",
		strings.ToUpper(result.SystemName),
	)

	fmt.Fprintf(
		&builder,
		"🛡️ Security: %.2f\n",
		result.SecurityStatus,
	)

	fmt.Fprintf(
		&builder,
		"🪐 Planets: %d\n",
		result.Planets,
	)

	fmt.Fprintf(
		&builder,
		"🚪 Stargates: %d\n",
		result.Stargates,
	)

	fmt.Fprintf(
		&builder,
		"🏭 Stations: %d\n\n",
		result.Stations,
	)

	builder.WriteString("📊 LAST HOUR\n")

	fmt.Fprintf(
		&builder,
		"🚀 Ship jumps: %d\n",
		result.ShipJumps,
	)

	fmt.Fprintf(
		&builder,
		"💀 Ship kills: %d\n",
		result.ShipKills,
	)

	fmt.Fprintf(
		&builder,
		"☠️ Pod kills: %d\n",
		result.PodKills,
	)

	fmt.Fprintf(
		&builder,
		"🤖 NPC kills: %d",
		result.NPCKills,
	)

	return builder.String()
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
