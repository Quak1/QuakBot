package main

import (
	"flag"
	"log"

	"github.com/Quak1/QuakBot/internal/riotapi"
)

var (
	championsURL string
	patch        string
	queuesURL    string
)

func main() {
	flag.StringVar(&championsURL, "c", "", "Champions URL")
	flag.StringVar(&patch, "p", "", "Patch version")
	flag.StringVar(&queuesURL, "q", "", "Queues URL")
	flag.Parse()

	if patch == "" && championsURL == "" {
		log.Panic("Either champions URL or Patch must be provided.")
	}
	if championsURL == "" {
		championsURL = riotapi.ChampionURL(patch)
	}
	if queuesURL == "" {
		queuesURL = riotapi.QueuesURL
	}

	log.Println("Downloading...")
	riotapi.DownloadConstants(championsURL, queuesURL)
	log.Println("Done")
}
