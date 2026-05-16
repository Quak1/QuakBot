package rpsgame

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type activeGame struct {
	userID string
	object string
}

var activeGames = map[string]activeGame{}

var rpsChoices = map[string]struct {
	description string
}{
	"rock": {
		description: "sedimentary, igneous, or perhaps even metamorphic",
	},
	"paper": {
		description: "versatile and iconic",
	},
	"scissors": {
		description: "careful ! sharp ! edges !!",
	},
}

func getOptions() []discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{}

	for k, o := range rpsChoices {
		options = append(options, discordgo.SelectMenuOption{
			Label:       firstToUpper(k),
			Value:       k,
			Description: o.description,
		})
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options
}

func handleAccept(s *discordgo.Session, i *discordgo.InteractionCreate) {
	gameID := strings.ReplaceAll(i.MessageComponentData().CustomID, "accept_button-", "")

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "What is your object of choice?",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							MenuType: discordgo.StringSelectMenu,
							CustomID: fmt.Sprintf("select_choice-%s", gameID),
							Options:  getOptions(),
						},
					},
				},
			},
		},
	})
}
