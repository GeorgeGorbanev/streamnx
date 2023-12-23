package main

import (
	"fmt"
	"log"
	"os"

	"github.com/GeorgeGorbanev/songshift/internal/songshift"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"

	"github.com/joho/godotenv"
)

type config struct {
	telegramToken string
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

	ss := songshift.NewSongshift(bot.Sender())
	bot.HandleText(ss.HandleText)
	defer bot.Stop()
	bot.Start()
}

func readConfig() (*config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &config{
		telegramToken: os.Getenv("TELEGRAM_TOKEN"),
	}, nil
}
