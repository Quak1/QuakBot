package observer

import (
	"fmt"
	"log"

	"github.com/Quak1/QuakBot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

// monitoredUserID -> registered in guildID -> notification to channelID
type UserStatus struct {
	status discordgo.Status
	guild  map[string]map[string]struct{}
}

var MonitoredUsers = map[string]*UserStatus{}

func GetObserverCommands() []*discordgo.ApplicationCommand {
	cmds := []*discordgo.ApplicationCommand{
		{
			Name:        "watch",
			Description: "Watch the online status of a user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Pick a user",
					Required:    true,
				},
			},
		},
	}

	return cmds
}

func GetObserverHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"watch": handleWatchCmd,
	}
}

func addMonitoredUser(userID, guildID, channelID string) {
	if MonitoredUsers[userID] == nil {
		MonitoredUsers[userID] = &UserStatus{
			guild: make(map[string]map[string]struct{}),
		}
	}
	if MonitoredUsers[userID].guild[guildID] == nil {
		MonitoredUsers[userID].guild[guildID] = make(map[string]struct{})
	}
	MonitoredUsers[userID].guild[guildID][channelID] = struct{}{}
}

func handleWatchCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	option := i.ApplicationCommandData().Options[0].Value.(string)

	userID := utils.GetUserID(i)
	ch, err := s.UserChannelCreate(userID)
	if err != nil {
		log.Println("Error creating channel:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something went wrong while sending the DM",
			},
		})
		return
	}

	addMonitoredUser(option, i.GuildID, ch.ID)

	_, err = s.ChannelMessageSend(ch.ID, fmt.Sprintf("Watching <@%s>", option))
	if err != nil {
		log.Println("Error sending DM:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "Failed to send you a DM",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "OK! Look at your DMs",
		},
	})
}

func HandlePresenceEvent(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	userID := p.User.ID
	if _, ok := MonitoredUsers[userID]; !ok {
		return
	}

	if channels, ok := MonitoredUsers[userID].guild[p.GuildID]; ok {
		for _, a := range p.Activities {
			log.Printf("%+v\n", a)
		}
		log.Println("Activities ----^")

		if MonitoredUsers[userID].status == p.Status {
			return
		} else {
			MonitoredUsers[userID].status = p.Status
		}

		for c := range channels {
			var msg string
			g, err := s.Guild(p.GuildID)
			if err != nil {
				msg = fmt.Sprintf("User <@%s> changed status to `%s`\n", p.User.ID, p.Status)
			} else {
				msg = fmt.Sprintf("User <@%s> changed status to `%s` in `%s`\n", p.User.ID, p.Status, g.Name)
			}

			_, err = s.ChannelMessageSend(c, msg)
			if err != nil {
				log.Println("[observer] Error sending DM:", err)
			}
		}
	}
}
