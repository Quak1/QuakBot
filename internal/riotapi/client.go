package riotapi

import (
	"net/http"
	"time"
)

type Client struct {
	Account   *AccountClient
	Spectator *SpectatorClient
	League    *LeagueClient
	Champions ChampionNameByID
	Queues    QueuesByID
}

func NewClient(region Region, apiKey string) (*Client, error) {
	baseClient := newBaseClient(region, apiKey, &http.Client{
		Timeout: 30 * time.Second,
	})

	var champions ChampionNameByID
	var queues QueuesByID
	if err := loadConstant(championsFilename, &champions); err != nil {
		return nil, err
	}
	if err := loadConstant(queuesFilename, &queues); err != nil {
		return nil, err
	}

	return &Client{
		Account:   newAccountClient(baseClient),
		Spectator: newSpectatorClient(baseClient),
		League:    newLeagueClient(baseClient),
		Champions: champions,
		Queues:    queues,
	}, nil
}
