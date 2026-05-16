package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var commands = []*discordgo.ApplicationCommand{
	{
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
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file", err)
	}

	token := os.Getenv("DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening sessions", err)
	}
	defer dg.Close()

	appID := dg.State.User.ID

	log.Println("Registering commands...")
	for _, c := range commands {
		_, err := dg.ApplicationCommandCreate(appID, "", c)
		if err != nil {
			log.Panicf("error creating command %v: %v", c.Name, err)
		}
	}

	cmds, err := dg.ApplicationCommands(appID, "")
	if err != nil {
		log.Panic("error retrieving commands: ", err)
	}

	fmt.Println("\nCommands:")
	for _, cmd := range cmds {
		fmt.Printf("- %s: %s\n", cmd.Name, cmd.Description)
	}
}
