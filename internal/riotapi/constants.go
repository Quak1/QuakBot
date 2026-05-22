package riotapi

import (
	"encoding/gob"
	"fmt"
	"os"
	"strconv"

	"github.com/Quak1/QuakBot/internal/utils"
)

const (
	championURLFormat = "https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json"
	QueuesURL         = "https://static.developer.riotgames.com/docs/lol/queues.json"
	championsFilename = "champions.gob"
	queuesFilename    = "queues.gob"
)

type queueDTO struct {
	QueueID     int    `json:"queueId"`
	Map         string `json:"map"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
}
type QueuesByID map[int]queueDTO

type ddChampions struct {
	Data map[string]struct {
		Key  string `json:"key"`
		Name string `json:"name"`
		ID   string `json:"id"`
	}
}
type ChampionNameByID map[int]string

func DownloadConstants(championURL, queuesURL string) error {
	var cs ddChampions
	var qs []queueDTO

	if err := utils.GetAndParseJSON(championURL, &cs); err != nil {
		return err
	}
	if err := utils.GetAndParseJSON(queuesURL, &qs); err != nil {
		return err
	}

	champions := make(ChampionNameByID)
	for _, c := range cs.Data {
		key, err := strconv.Atoi(c.Key)
		if err != nil {
			return err
		}
		champions[key] = c.ID
	}

	queues := make(QueuesByID)
	for _, q := range qs {
		queues[q.QueueID] = q
	}

	if err := saveConstant(championsFilename, &champions); err != nil {
		return err
	}
	if err := saveConstant(queuesFilename, &queues); err != nil {
		return err
	}

	return nil
}

func ChampionURL(patch string) string {
	return fmt.Sprintf(championURLFormat, patch)
}

func saveConstant(filename string, data any) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(data); err != nil {
		return err
	}

	return nil
}

func loadConstant(filename string, data any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	if err := dec.Decode(data); err != nil {
		return err
	}

	return nil
}
