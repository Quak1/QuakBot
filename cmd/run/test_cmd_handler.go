package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

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
