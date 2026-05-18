package main

import (
	"github.com/bwmarrin/discordgo"
)

var testCommand = &discordgo.ApplicationCommand{
	Name:        "test",
	Description: "Basic command",
	Type:        discordgo.ChatApplicationCommand,
	IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
		discordgo.ApplicationIntegrationGuildInstall,
		discordgo.ApplicationIntegrationUserInstall,
	},
	Contexts: &[]discordgo.InteractionContextType{
		discordgo.InteractionContextGuild,
		discordgo.InteractionContextBotDM,
		discordgo.InteractionContextPrivateChannel,
	},
}
