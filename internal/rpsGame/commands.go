package rpsgame

import (
	"github.com/Quak1/QuakBot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func GetRPSCommands() []*discordgo.ApplicationCommand {
	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "challenge",
			Description: "Challenge to a match of rock paper scissors",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "Pick your item",
					Required:    true,
					Choices:     createChallengeCmdChoices(),
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

	return cmds
}

func createChallengeCmdChoices() []*discordgo.ApplicationCommandOptionChoice {
	commandChoices := []*discordgo.ApplicationCommandOptionChoice{}

	for k := range rpsObjects {
		commandChoices = append(commandChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  utils.FirstLetterToUpper(k),
			Value: k,
		})
	}

	return commandChoices
}
