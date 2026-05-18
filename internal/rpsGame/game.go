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
	matchups    map[string]string
}{
	"rock": {
		description: "sedimentary, igneous, or perhaps even metamorphic",
		matchups: map[string]string{
			"scissors": "crushes",
		},
	},
	"paper": {
		description: "versatile and iconic",
		matchups: map[string]string{
			"rock": "covers",
		},
	},
	"scissors": {
		description: "careful ! sharp ! edges !!",
		matchups: map[string]string{
			"paper": "cuts",
		},
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
	s.ChannelMessageDelete(i.ChannelID, i.Message.ID)

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

func handleItemChoice(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	gameID := strings.ReplaceAll(data.CustomID, "select_choice-", "")

	if game, ok := activeGames[gameID]; ok {
		userID := getUserID(i)
		object := data.Values[0]
		result := getResult(game, activeGame{userID: userID, object: object})

		delete(activeGames, gameID)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: "Nice choice!",
			},
		})

		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{Content: result})
	}
}

func getResult(p1, p2 activeGame) string {
	if msg, ok := rpsChoices[p1.object].matchups[p2.object]; ok {
		return formatResult(p1, p2, msg)
	} else if msg, ok := rpsChoices[p2.object].matchups[p1.object]; ok {
		return formatResult(p2, p1, msg)
	} else {
		return formatResult(p1, p2, "tie")
	}
}

func formatResult(win, lose activeGame, msg string) string {
	if msg == "tie" {
		return fmt.Sprintf("<@%s> and <@%s> draw with **%s**", win.userID, lose.userID, win.object)
	}

	return fmt.Sprintf("<@%s>'s **%s** %s <@%s>'s **%s**", win.userID, win.object, msg, lose.userID, lose.object)
}
