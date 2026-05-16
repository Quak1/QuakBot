package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmdName := i.ApplicationCommandData().Name
	if cmdName == "test" {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Test command ran successfully!",
			},
		})
		if err != nil {
			log.Println(err)
		}
	} else {
		log.Println("Received unknown command: ", cmdName)
	}
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

	dg.AddHandler(commandHandler)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
