package main

import (
	"fmt"
	"log"
	"maps"
	"os"
	"os/signal"
	"strings"
	"syscall"

	rpsgame "github.com/Quak1/QuakBot/internal/rpsGame"
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

func testCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	commands = append(commands, rpsgame.GetRPSCommands()...)
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

	handlers := map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"test": testCommandHandler,
	}
	maps.Copy(handlers, rpsgame.GetRPSCommandHandlers())
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		var name string
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			name = i.ApplicationCommandData().Name
		case discordgo.InteractionMessageComponent:
			name = i.MessageComponentData().CustomID
			if j := strings.IndexByte(name, '-'); j != -1 {
				name = name[:j]
			}
		}

		fmt.Println("command received:", name)

		if h, ok := handlers[name]; ok {
			h(s, i)
		}
	})

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
