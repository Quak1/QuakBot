package riotbot

import "github.com/bwmarrin/discordgo"

func GetCommands() []*discordgo.ApplicationCommand {
	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "livegame",
			Description: "Get live game information",
			Options: []*discordgo.ApplicationCommandOption{
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
