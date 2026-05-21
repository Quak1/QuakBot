package riotapi

import (
	"fmt"
)

type LeagueClient struct {
	c *baseClient
}

type LeagueEntryDTO struct {
	LeagueID     string         `json:"leagueId"`
	PUUID        string         `json:"puuid"`
	QueueType    string         `json:"queueType"`
	Tier         string         `json:"tier"`
	Rank         string         `json:"rank"`
	LeaguePoints int            `json:"leaguePoints"`
	Wins         int            `json:"wins"`
	Losses       int            `json:"losses"`
	HotStreak    bool           `json:"hotStreak"`
	Veteran      bool           `json:"veteran"`
	FreshBlood   bool           `json:"freshBlood"`
	Inactive     bool           `json:"inactive"`
	MiniSeries   *MiniSeriesDTO `json:"miniSeries"`
}

type MiniSeriesDTO struct {
	Losses   int    `json:"losses"`
	Progress string `json:"progress"`
	Target   int    `json:"target"`
	Wins     int    `json:"wins"`
}

const leagueURL = "lol/league/v4"

func newLeagueClient(c *baseClient) *LeagueClient {
	return &LeagueClient{
		c: c,
	}
}

func (lc *LeagueClient) GetLeaguesByPUUID(puuid string) ([]*LeagueEntryDTO, error) {
	data := []*LeagueEntryDTO{}
	if err := lc.c.DoGet(fmt.Sprintf("%s/entries/by-puuid/%s", leagueURL, puuid), &data); err != nil {
		return nil, err
	}

	return data, nil
}
