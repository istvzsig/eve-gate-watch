package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/istvzsig/eve-gate-watch/internal/intel"
	"github.com/istvzsig/eve-gate-watch/internal/utils"
)

const (
	PingCommand  = "!ping"
	DScanCommand = "!dscan"
	PVPCommand   = "!pvp"
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

	// pvp
	if msgContent == PVPCommand || strings.HasPrefix(msgContent, PVPCommand+" ") {
		dispatchPVPCommand(session, message, intelService)
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

	messageText := utils.FormatDScan(result)

	_, err = session.ChannelMessageSend(
		message.ChannelID,
		messageText,
	)
	if err != nil {
		log.Printf("send dscan result: %v", err)
	}
}

func dispatchPVPCommand(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
	intelService *intel.Service,
) {
	systemName := strings.TrimSpace(
		strings.TrimPrefix(message.Content, PVPCommand),
	)

	result, err := intelService.RecentPVP(systemName)
	if err != nil {
		log.Printf("pvp: %v", err)

		_, sendErr := session.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprintf(
				"System `%s` not found.",
				systemName,
			),
		)
		if sendErr != nil {
			log.Printf("send pvp error: %v", sendErr)
		}

		return
	}

	messageText := utils.FormatPVP(result)

	_, err = session.ChannelMessageSend(
		message.ChannelID,
		messageText,
	)
	if err != nil {
		log.Printf("send pvp result: %v", err)
	}
}
