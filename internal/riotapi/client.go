package riotapi

import (
	"net/http"
	"time"
)

type Client struct {
	Account *AccountClient
}

func NewClient(region Region, apiKey string) *Client {
	baseClient := newBaseClient(region, apiKey, &http.Client{
		Timeout: 30 * time.Second,
	})

	return &Client{
		Account: NewAccountClient(baseClient),
	}
}
