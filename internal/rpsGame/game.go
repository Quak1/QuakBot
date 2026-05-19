package rpsgame

import (
	"fmt"
)

type activeGame struct {
	userID string
	object string
}

var activeGames = map[string]activeGame{}

var rpsObjects = map[string]struct {
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

func getResult(p1, p2 activeGame) string {
	if msg, ok := rpsObjects[p1.object].matchups[p2.object]; ok {
		return formatResult(p1, p2, msg)
	} else if msg, ok := rpsObjects[p2.object].matchups[p1.object]; ok {
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
