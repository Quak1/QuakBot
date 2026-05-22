package riotapi

import (
	"fmt"
)

type SpectatorClient struct {
	c *baseClient
}

type CurrentGameInfo struct {
	GameID            int64                     `json:"gameId"`
	GameType          string                    `json:"gameType"`
	GameStartTime     int64                     `json:"gameStartTime"`
	MapID             int64                     `json:"mapId"`
	GameLength        int64                     `json:"gameLength"`
	PlatformID        string                    `json:"platformId"`
	GameMode          string                    `json:"gameMode"`
	BannedChampions   []*BannedChampions        `json:"bannedChampions"`
	GameQueueConfigID int64                     `json:"gameQueueConfigId"`
	Observers         *Observer                 `json:"observers"`
	Participants      []*CurrentGameParticipant `json:"participants"`
}

type BannedChampions struct {
	PickTurn   int   `json:"pickTurn"`
	ChampionID int64 `json:"championId"`
	TeamID     int64 `json:"teamId"`
}

type CurrentGameParticipant struct {
	ChampionID               int64                      `json:"championId"`
	Perks                    *Perks                     `json:"perks"`
	ProfileIconID            int64                      `json:"profileIconId"`
	Bot                      bool                       `json:"bot"`
	TeamID                   int64                      `json:"teamId"`
	RiotID                   string                     `json:"riotId"`
	PUUID                    string                     `json:"puuid"`
	Spell1ID                 int64                      `json:"spell1Id"`
	Spell2ID                 int64                      `json:"spell2Id"`
	GameCustomizationObjects []*GameCustomizationObject `json:"gameCustomizationObjects"`
}
type Perks struct {
	PerkIds      []int64 `json:"perkIds"`
	PerkStyle    int64   `json:"perkStyle"`
	PerkSubStyle int64   `json:"perkSubStyle"`
}

type GameCustomizationObject struct {
	Category string `json:"string"`
	Content  string `json:"content"`
}

type Observer struct {
	EncryptionKey string `json:"encryptionKey"`
}

const spectatorURL = "lol/spectator/v5/active-games/by-summoner"

func newSpectatorClient(c *baseClient) *SpectatorClient {
	return &SpectatorClient{
		c: c,
	}
}

func (sc *SpectatorClient) GetGameByPUUID(puuid string) (*CurrentGameInfo, error) {
	data := &CurrentGameInfo{}
	if err := sc.c.DoGet(fmt.Sprintf("%s/%s", spectatorURL, puuid), data); err != nil {
		return nil, err
	}

	return data, nil
}
