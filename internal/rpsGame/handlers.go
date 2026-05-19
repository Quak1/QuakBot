package rpsgame

import (
	"fmt"
	"log"
	"math/rand/v2"
	"strings"

	"github.com/Quak1/QuakBot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

const (
	AcceptChallenge = "accept_button-"
	SelectChoice    = "select_choice-"
)

func GetRPSHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"challenge":     handleChallengeCmd,
		"accept_button": handleAcceptChallengeInt,
		"select_choice": handleObjectChoiceInt,
	}
}

func handleChallengeCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := utils.GetUserID(i)
	option := i.ApplicationCommandData().Options[0].Value

	activeGames[i.ID] = activeGame{
		userID: userID,
		object: option.(string),
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Rock papers scissors challenge from <@%s>", userID),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: fmt.Sprintf("%s%s", AcceptChallenge, i.ID),
							Label:    "Accept",
							Style:    discordgo.PrimaryButton,
						},
					},
				},
			},
		},
	})

	if err != nil {
		log.Println(err)
	}
}

func getOptions() []discordgo.SelectMenuOption {
	options := []discordgo.SelectMenuOption{}

	for k, o := range rpsObjects {
		options = append(options, discordgo.SelectMenuOption{
			Label:       utils.FirstLetterToUpper(k),
			Value:       k,
			Description: o.description,
		})
	}

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options
}

func handleAcceptChallengeInt(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.ChannelMessageDelete(i.ChannelID, i.Message.ID)

	gameID := strings.ReplaceAll(i.MessageComponentData().CustomID, AcceptChallenge, "")

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
							CustomID: fmt.Sprintf("%s%s", SelectChoice, gameID),
							Options:  getOptions(),
						},
					},
				},
			},
		},
	})
}

func handleObjectChoiceInt(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.MessageComponentData()
	gameID := strings.ReplaceAll(data.CustomID, SelectChoice, "")

	if game, ok := activeGames[gameID]; ok {
		userID := utils.GetUserID(i)
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
