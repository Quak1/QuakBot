package riotbot

import (
	"github.com/Quak1/QuakBot/internal/riotapi"
	"github.com/bwmarrin/discordgo"
)

func GetCommands() []*discordgo.ApplicationCommand {
	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "livegame",
			Description: "Get live game information",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "region",
					Description: "Region",
					Required:    true,
					Choices:     getRegionChoices(),
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "Game name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tag",
					Description: "Tag line",
					Required:    true,
				},
			},
		},
	}

	return cmds
}

func getRegionChoices() []*discordgo.ApplicationCommandOptionChoice {
	choices := []*discordgo.ApplicationCommandOptionChoice{}
	for _, r := range riotapi.Regions {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  string(r),
			Value: r,
		})
	}

	return choices
}
