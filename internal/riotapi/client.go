package riotapi

import (
	"net/http"
	"time"
)

type Client struct {
	Account   *AccountClient
	Spectator *SpectatorClient
	League    *LeagueClient
}

func NewClient(region Region, apiKey string) *Client {
	baseClient := newBaseClient(region, apiKey, &http.Client{
		Timeout: 30 * time.Second,
	})

	return &Client{
		Account:   NewAccountClient(baseClient),
		Spectator: newSpectatorClient(baseClient),
		League:    newLeagueClient(baseClient),
	}
}
