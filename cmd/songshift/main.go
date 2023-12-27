package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GeorgeGorbanev/songshift/internal/songshift"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"

	"github.com/joho/godotenv"
)

type config struct {
	telegramToken       string
	spotifyClientID     string
	spotifyClientSecret string
}

func main() {
	cfg, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	bot, err := telegram.NewBot(cfg.telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	spotifyClient := makeSpotifyClient(cfg)
	ss := songshift.NewSongshift(spotifyClient, bot.Sender())
	bot.HandleText(ss.HandleText)
	defer bot.Stop()
	bot.Start()
}

func readConfig() (*config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &config{
		telegramToken:       os.Getenv("TELEGRAM_TOKEN"),
		spotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}, nil
}

func makeSpotifyClient(cfg *config) *spotify.Client {
	return spotify.NewClient(&spotify.Credentials{
		ClientID:     cfg.spotifyClientID,
		ClientSecret: cfg.spotifyClientSecret,
	})
}
