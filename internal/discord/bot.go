package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/istvzsig/eve-gate-watch/internal/intel"
	"github.com/istvzsig/eve-gate-watch/internal/systems"
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
	return session.Open()
}

func CloseSession(session *discordgo.Session) error {
	return session.Close()
}

// RegisterPrefixCommands keeps the existing prefix commands and also
// registers the interaction dispatcher used by slash commands.
func RegisterPrefixCommands(
	session *discordgo.Session,
	intelService *intel.Service,
	systemsService *systems.Service,
) {
	// Prefix commands:
	// !ping
	// !dscan <system>
	// !pvp <system>
	session.AddHandler(func(
		session *discordgo.Session,
		message *discordgo.MessageCreate,
	) {
		DispatchPrefixCommands(session, message, intelService)
	})

	// Slash commands + autocomplete.
	session.AddHandler(func(
		session *discordgo.Session,
		interaction *discordgo.InteractionCreate,
	) {
		if interaction.Type != discordgo.InteractionApplicationCommand &&
			interaction.Type != discordgo.InteractionApplicationCommandAutocomplete {
			return
		}

		DispatchSlashCommands(
			session,
			interaction,
			intelService,
			systemsService,
		)
	})
}

func RegisterSlashCommands(session *discordgo.Session) error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "dscan",
			Description: "Show system intelligence",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "system",
					Description:  "EVE system",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "pvp",
			Description: "Show recent PvP activity",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "system",
					Description:  "EVE system",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	for _, command := range commands {
		_, err := session.ApplicationCommandCreate(
			session.State.User.ID,
			"",
			command,
		)
		if err != nil {
			return fmt.Errorf(
				"register /%s: %w",
				command.Name,
				err,
			)
		}
	}

	return nil
}

// ------------------------------------------------------------
// Prefix commands
// ------------------------------------------------------------

func DispatchPrefixCommands(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
	intelService *intel.Service,
) {
	if message.Author == nil || message.Author.Bot {
		return
	}

	content := strings.TrimSpace(message.Content)

	switch {
	case content == PingCommand:
		dispatchPingCommand(session, message)

	case strings.HasPrefix(content, DScanCommand):
		dispatchDScanCommand(session, message, intelService)

	case strings.HasPrefix(content, PVPCommand):
		dispatchPVPCommand(session, message, intelService)
	}
}

func dispatchPingCommand(
	session *discordgo.Session,
	message *discordgo.MessageCreate,
) {
	_, err := session.ChannelMessageSend(
		message.ChannelID,
		"pong",
	)
	if err != nil {
		log.Printf("send ping response: %v", err)
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

	if systemName == "" {
		_, err := session.ChannelMessageSend(
			message.ChannelID,
			"Usage: `!dscan <system>`",
		)
		if err != nil {
			log.Printf("send dscan usage: %v", err)
		}
		return
	}

	result, err := intelService.DScan(systemName)
	if err != nil {
		log.Printf("dscan %s: %v", systemName, err)

		_, sendErr := session.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprintf(
				"Could not load intelligence for `%s`.",
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

	if systemName == "" {
		_, err := session.ChannelMessageSend(
			message.ChannelID,
			"Usage: `!pvp <system>`",
		)
		if err != nil {
			log.Printf("send pvp usage: %v", err)
		}
		return
	}

	result, err := intelService.RecentPVP(systemName)
	if err != nil {
		log.Printf("pvp %s: %v", systemName, err)

		_, sendErr := session.ChannelMessageSend(
			message.ChannelID,
			fmt.Sprintf(
				"Could not load PvP data for `%s`.",
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

// ------------------------------------------------------------
// Slash commands
// ------------------------------------------------------------

func DispatchSlashCommands(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	intelService *intel.Service,
	systemsService *systems.Service,
) {
	data := interaction.ApplicationCommandData()

	if interaction.Type == discordgo.InteractionApplicationCommandAutocomplete {
		dispatchSystemAutocomplete(
			session,
			interaction,
			systemsService,
		)
		return
	}

	switch data.Name {
	case "dscan":
		dispatchSlashDScan(
			session,
			interaction,
			intelService,
		)

	case "pvp":
		dispatchSlashPVP(
			session,
			interaction,
			intelService,
		)
	}
}

func dispatchSystemAutocomplete(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	systemsService *systems.Service,
) {
	data := interaction.ApplicationCommandData()

	if len(data.Options) == 0 {
		return
	}

	query := data.Options[0].StringValue()

	results := systemsService.Search(query)

	choices := make(
		[]*discordgo.ApplicationCommandOptionChoice,
		0,
		len(results),
	)

	for _, system := range results {
		choices = append(
			choices,
			&discordgo.ApplicationCommandOptionChoice{
				Name:  system,
				Value: system,
			},
		)

		if len(choices) >= 25 {
			break
		}
	}

	err := session.InteractionRespond(
		interaction.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		},
	)

	if err != nil {
		log.Printf("system autocomplete: %v", err)
	}
}

func dispatchSlashDScan(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	intelService *intel.Service,
) {
	data := interaction.ApplicationCommandData()

	if len(data.Options) == 0 {
		return
	}

	systemName := data.Options[0].StringValue()

	// Acknowledge the interaction immediately.
	err := session.InteractionRespond(
		interaction.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		},
	)
	if err != nil {
		log.Printf("defer dscan interaction: %v", err)
		return
	}

	result, err := intelService.DScan(systemName)
	if err != nil {
		log.Printf("slash dscan %s: %v", systemName, err)

		message := fmt.Sprintf(
			"Could not load intelligence for `%s`.",
			systemName,
		)

		_, editErr := session.InteractionResponseEdit(
			interaction.Interaction,
			&discordgo.WebhookEdit{
				Content: &message,
			},
		)
		if editErr != nil {
			log.Printf("edit dscan error response: %v", editErr)
		}

		return
	}

	messageText := utils.FormatDScan(result)

	_, err = session.InteractionResponseEdit(
		interaction.Interaction,
		&discordgo.WebhookEdit{
			Content: &messageText,
		},
	)
	if err != nil {
		log.Printf("send dscan result: %v", err)
	}
}

func dispatchSlashPVP(
	session *discordgo.Session,
	interaction *discordgo.InteractionCreate,
	intelService *intel.Service,
) {
	data := interaction.ApplicationCommandData()

	if len(data.Options) == 0 {
		return
	}

	systemName := data.Options[0].StringValue()

	// IMPORTANT:
	// zKill/ESI can take longer than Discord's initial interaction
	// response window, so acknowledge the interaction first.
	err := session.InteractionRespond(
		interaction.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		},
	)
	if err != nil {
		log.Printf("defer pvp interaction: %v", err)
		return
	}

	result, err := intelService.RecentPVP(systemName)
	if err != nil {
		log.Printf("slash pvp %s: %v", systemName, err)

		message := fmt.Sprintf(
			"Could not load PvP data for `%s`.",
			systemName,
		)

		_, editErr := session.InteractionResponseEdit(
			interaction.Interaction,
			&discordgo.WebhookEdit{
				Content: &message,
			},
		)
		if editErr != nil {
			log.Printf("edit pvp error response: %v", editErr)
		}

		return
	}

	messageText := utils.FormatPVP(result)

	_, err = session.InteractionResponseEdit(
		interaction.Interaction,
		&discordgo.WebhookEdit{
			Content: &messageText,
		},
	)
	if err != nil {
		log.Printf("send pvp result: %v", err)
	}
}
