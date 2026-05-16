package rpsgame

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func createCommandChoices() []*discordgo.ApplicationCommandOptionChoice {
	commandChoices := []*discordgo.ApplicationCommandOptionChoice{}

	for k := range rpsChoices {
		commandChoices = append(commandChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  firstToUpper(k),
			Value: k,
		})
	}

	return commandChoices
}

func firstToUpper(s string) string {
	if len(s) != 0 && s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}

	return s
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "challenge",
			Description: "Challenge to a match of rock paper scissors",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "Pick your item",
					Required:    true,
					Choices:     createCommandChoices(),
				},
			},
			Type: discordgo.ChatApplicationCommand,
			IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
				discordgo.ApplicationIntegrationGuildInstall,
				discordgo.ApplicationIntegrationUserInstall,
			},
			Contexts: &[]discordgo.InteractionContextType{
				discordgo.InteractionContextGuild,
				discordgo.InteractionContextPrivateChannel,
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"challenge": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.User
			if user == nil {
				user = i.Member.User
			}
			userID := user.ID
			option := i.ApplicationCommandData().Options[0].Value

			activeGames[i.ID] = activeGame{
				userID: userID,
				object: option.(string),
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Rock papers scissors challenge from <@%s> : %s", userID, option),
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									CustomID: fmt.Sprintf("accept_button-%s", i.ID),
									Label:    "Accept",
									Style:    discordgo.PrimaryButton,
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Println(err)
			}
		},
		"accept_button": handleAccept,
	}
)

func GetRPSCommands() []*discordgo.ApplicationCommand {
	return commands
}

func GetRPSCommandHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate) {
	return commandHandlers
}
