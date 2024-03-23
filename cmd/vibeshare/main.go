package main

import (
	"log/slog"
	"os"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/converter"
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
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{AddSource: true},
			),
		),
	)

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			slog.Error("failed to load .env file", slog.String("error", err.Error()))
			return
		}
	}
	cfg := newConfig()

	bot, err := telegram.NewBot(cfg.telegramToken)
	if err != nil {
		slog.Error("failed to create bot", slog.String("error", err.Error()))
		return
	}

	vs := newVibeshare(cfg, bot.Sender())
	bot.HandleText(vs.HandleText)
	defer bot.Stop()

	slog.Info("Bot started")
	bot.Start()
}

func newConfig() *config {
	return &config{
		telegramToken:       os.Getenv("TELEGRAM_TOKEN"),
		spotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
	}
}

func newVibeshare(cfg *config, ts telegram.Sender) *vibeshare.Vibeshare {
	return vibeshare.NewVibeshare(&vibeshare.Input{
		Converter: converter.NewConverter(&converter.Input{
			SpotifyClient: spotify.NewClient(&spotify.Credentials{
				ClientID:     cfg.spotifyClientID,
				ClientSecret: cfg.spotifyClientSecret,
			}),
			YandexClient: ymusic.NewClient(),
		}),
		TelegramSender: ts,
	})
}
