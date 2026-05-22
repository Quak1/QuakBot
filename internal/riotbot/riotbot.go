package riotbot

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Quak1/QuakBot/internal/riotapi"
	"github.com/bwmarrin/discordgo"
)

type player struct {
	puuid          string
	name           string
	championNameID string
	solo           *queueRank
	flex           *queueRank
}

type queueRank struct {
	league   string
	division string
	wins     int
	loses    int
}

type gameTeams struct {
	blue []*player
	red  []*player
}

type result struct {
	player *riotapi.CurrentGameParticipant
	ranks  []*riotapi.LeagueEntryDTO
	err    error
}

var riotClients = map[string]*riotapi.Client{}

const (
	ddImgFormat = "https://ddragon.leagueoflegends.com/cdn/img/champion/tiles/%s_0.jpg"
)

var (
	RIOT_TOKEN    = os.Getenv("RIOT_TOKEN")
	leagueToEmoji = map[string]string{
		"IRON":        "poop",
		"BRONZE":      "brown_circle",
		"SILVER":      "white_circle",
		"GOLD":        "yellow_circle",
		"PLATINUM":    "blue_circle",
		"EMERALD":     "green_circle",
		"DIAMOND":     "gem",
		"MASTER":      "purple_circle",
		"GRANDMASTER": "diamonds",
		"CHALLENGER":  "rosette",
		"":            "x",
		"STREAMER":    "clown",
	}
	divisionToEmoji = map[string]string{
		"I":   "number_1",
		"II":  "number_2",
		"III": "number_3",
		"IV":  "number_4",
		"":    "number_0",
	}
)

func getGamePlayerInfo(rc *riotapi.Client, game *riotapi.CurrentGameInfo) (*gameTeams, error) {
	var wg sync.WaitGroup
	ch := make(chan result, len(game.Participants))

	for _, p := range game.Participants {
		wg.Add(1)
		go func(p *riotapi.CurrentGameParticipant) {
			defer wg.Done()

			if p.PUUID == "" {
				ch <- result{player: p}
				return
			}

			rank, err := rc.League.GetLeaguesByPUUID(p.PUUID)
			if err != nil {
				ch <- result{err: err}
				return
			}

			ch <- result{player: p, ranks: rank}
		}(p)
	}

	wg.Wait()
	close(ch)

	teams := gameTeams{
		blue: []*player{},
		red:  []*player{},
	}

	for res := range ch {
		if res.err != nil {
			return nil, res.err
		}

		participant := player{
			puuid:          res.player.PUUID,
			name:           res.player.RiotID,
			championNameID: rc.Champions[int(res.player.ChampionID)],
			solo:           &queueRank{},
		}

		if res.player.PUUID == "" {
			participant.solo = &queueRank{
				league:   "STREAMER",
				division: "I",
			}
		}

		for _, r := range res.ranks {
			rank := &queueRank{
				league:   r.Tier,
				division: r.Rank,
				wins:     r.Wins,
				loses:    r.Losses,
			}
			switch r.QueueType {
			case "RANKED_SOLO_5x5":
				participant.solo = rank
			case "RANKED_FLEX_SR":
				participant.flex = rank
			}
		}

		switch res.player.TeamID {
		case 100:
			teams.blue = append(teams.blue, &participant)
		case 200:
			teams.red = append(teams.red, &participant)
		default:
			return nil, fmt.Errorf("Invalid team id: %d", res.player.TeamID)
		}
	}

	return &teams, nil
}

func parseTeam(team []*player, color int) discordgo.Container {
	components := []discordgo.MessageComponent{}

	for _, p := range team {
		components = append(components, parsePlayer(p))
	}

	return discordgo.Container{
		AccentColor: &color,
		Components:  components,
	}
}

func parsePlayer(p *player) discordgo.Section {
	strFormat := "   :%s: :%s:   `%03d/%03d %02d%%`"
	name := fmt.Sprintf("### %s\n", p.name)

	var soloWr int
	if p.solo.wins == 0 && p.solo.loses == 0 {
		soloWr = 0
	} else {
		soloWr = p.solo.wins * 100 / (p.solo.wins + p.solo.loses)
	}
	soloStr := "`S:`" + fmt.Sprintf(strFormat, leagueToEmoji[p.solo.league], divisionToEmoji[p.solo.division], p.solo.wins, p.solo.loses, soloWr)

	flexStr := ""
	if p.flex != nil {
		flexWr := p.flex.wins * 100 / (p.flex.wins + p.flex.loses)
		flexStr = "`F:`" + fmt.Sprintf(strFormat, leagueToEmoji[p.flex.league], divisionToEmoji[p.flex.division], p.flex.wins, p.flex.loses, flexWr)
	}

	if p.championNameID == "Fiddlesticks" {
		p.championNameID = "FiddleSticks"
	}

	return discordgo.Section{
		Accessory: discordgo.Thumbnail{
			Media: discordgo.UnfurledMediaItem{
				URL: fmt.Sprintf(ddImgFormat, strings.ReplaceAll(p.championNameID, " ", "")),
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{
				Content: name + soloStr + "\n" + flexStr,
			},
		},
	}
}
