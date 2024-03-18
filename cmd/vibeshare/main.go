package main

import (
	"log"
	"os"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/ymusic"

	"github.com/joho/godotenv"
)

type config struct {
	telegramToken       string
	spotifyClientID     string
	spotifyClientSecret string
}

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}

	cfg := config{
		telegramToken:       os.Getenv("TELEGRAM_TOKEN"),
		spotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}

	bot, err := telegram.NewBot(cfg.telegramToken)
	if err != nil {
		log.Fatal(err)
	}

	spotifyClient := spotify.NewClient(&spotify.Credentials{
		ClientID:     cfg.spotifyClientID,
		ClientSecret: cfg.spotifyClientSecret,
	})
	ts := vibeshare.NewVibeshare(&vibeshare.Input{
		SpotifyClient:  spotifyClient,
		TelegramSender: bot.Sender(),
		YmusicClient:   ymusic.NewClient(),
	})
	bot.HandleText(ts.HandleText)
	defer bot.Stop()

	bot.Start()
}
