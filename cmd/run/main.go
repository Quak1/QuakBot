package main

import (
	"flag"
	"log"
	"maps"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Quak1/QuakBot/internal/observer"
	rpsgame "github.com/Quak1/QuakBot/internal/rpsGame"
	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Discord bot token")
	flag.Parse()

	if Token == "" {
		Token = os.Getenv("DISCORD_TOKEN")
	}
}

func main() {
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	s.Identify.Intents |= discordgo.IntentGuildPresences

	handlers := map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		"test": testCommandHandler,
	}

	maps.Copy(handlers, rpsgame.GetRPSHandlers())
	maps.Copy(handlers, observer.GetObserverHandlers())
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

		if h, ok := handlers[name]; ok {
			h(s, i)
		}
	})

	s.AddHandler(observer.HandlePresenceEvent)

	err = s.Open()
	if err != nil {
		log.Fatal("error opening sessions", err)
	}
	defer s.Close()

	log.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
