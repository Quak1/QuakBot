package riotapi

import (
	"fmt"
)

type AccountClient struct {
	c *baseClient
}

type AccountDto struct {
	PUUID    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

const accountURL = "riot/account/v1/accounts"

func newAccountClient(c *baseClient) *AccountClient {
	return &AccountClient{
		c: c,
	}
}

func (ac *AccountClient) GetByRiotID(gameName, tagLine string) (*AccountDto, error) {
	data := &AccountDto{}
	if err := ac.c.DoAreaGet(fmt.Sprintf("%s/by-riot-id/%s/%s", accountURL, gameName, tagLine), data); err != nil {
		return nil, err
	}

	return data, nil
}
