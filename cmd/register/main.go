package main

import (
	"flag"
	"log"
	"os"

	rpsgame "github.com/Quak1/QuakBot/internal/rpsGame"
	"github.com/bwmarrin/discordgo"
)

var (
	Token             string
	DeleteOldCommands bool
)

func init() {
	flag.StringVar(&Token, "t", "", "Discord bot token")
	flag.BoolVar(&DeleteOldCommands, "d", false, "Delete all commands before registering")
	flag.Parse()

	if Token == "" {
		Token = os.Getenv("DISCORD_TOKEN")
	}
}

func main() {
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	err = s.Open()
	if err != nil {
		log.Fatal("error opening sessions", err)
	}
	defer s.Close()

	if DeleteOldCommands {
		deleteCommands(s, "")
	}

	var commands = []*discordgo.ApplicationCommand{testCommand}
	commands = append(commands, rpsgame.GetRPSCommands()...)

	registerCommands(s, "", commands)
}

func deleteCommands(s *discordgo.Session, guildID string) {
	appID := s.State.User.ID

	cmds, err := s.ApplicationCommands(appID, guildID)
	if err != nil {
		log.Fatal("Error fetching commands:", err)
	}

	for _, cmd := range cmds {
		err := s.ApplicationCommandDelete(appID, guildID, cmd.ID)
		if err != nil {
			log.Printf("Error deleting command %s: %v", cmd.Name, err)
		}
	}
}

func registerCommands(s *discordgo.Session, guildID string, cmds []*discordgo.ApplicationCommand) {
	log.Println("Registering commands...")

	for _, cmd := range cmds {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd)
		if err != nil {
			log.Panicf("Error creating command %v: %v", cmd.Name, err)
		}
	}
}
