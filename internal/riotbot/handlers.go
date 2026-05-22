package riotbot

import (
	"fmt"
	"log"
	"strings"

	"github.com/Quak1/QuakBot/internal/riotapi"
	"github.com/Quak1/QuakBot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func GetHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"livegame": handleLivegameCmd,
	}
}

func handleLivegameCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	region := data.Options[0].Value.(string)
	name := data.Options[1].Value.(string)
	tag := data.Options[2].Value.(string)
	log.Printf("Info request for: %s#%s @ %s", name, tag, region)

	if _, ok := riotClients[region]; !ok {
		rc, err := riotapi.NewClient(riotapi.Region(region), RIOT_TOKEN)
		if err != nil {
			utils.InteractionErrorResponse(s, i, err, "There was an error. Try again later!")
			return
		}

		riotClients[region] = rc
	}

	rc := riotClients[region]

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Loading...",
		},
	})

	user, err := rc.Account.GetByRiotID(name, tag)
	if err != nil {
		utils.InteractionFollowUpErrorResponse(s, i, err, "Error getting account information. Try again later!")
		return
	}

	game, err := rc.Spectator.GetGameByPUUID(user.PUUID)
	if err != nil {
		utils.InteractionFollowUpErrorResponse(s, i, err, "Error getting game information. Try again later!")
		return
	}

	teams, err := getGamePlayerInfo(rc, game)
	if err != nil {
		utils.InteractionFollowUpErrorResponse(s, i, err, "Error getting player information. Try again later!")
		return
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Flags: discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{
				Content: fmt.Sprintf("## %s", strings.Replace(rc.Queues[int(game.GameQueueConfigID)].Description, " games", "", 1)),
			},
			parseTeam(teams.blue, 0x0000FF),
			parseTeam(teams.red, 0xFF0000),
		},
	})
	if err != nil {
		log.Println(err)
	}
}
