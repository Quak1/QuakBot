package utils

import "github.com/bwmarrin/discordgo"

func FirstLetterToUpper(s string) string {
	if len(s) != 0 && s[0] >= 'a' && s[0] <= 'z' {
		return string(s[0]-32) + s[1:]
	}

	return s
}

func GetUserID(i *discordgo.InteractionCreate) string {
	user := i.User // only available on DMs
	if user == nil {
		user = i.Member.User
	}
	return user.ID
}
