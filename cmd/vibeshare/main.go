package main

import (
	"log/slog"
	"os"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/youtube"

	"github.com/joho/godotenv"
)

type config struct {
	telegramToken       string
	spotifyClientID     string
	spotifyClientSecret string
	youtubeAPIKey       string
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
			slog.Error("failed to load .env file", slog.Any("error", err))
			return
		}
	}
	cfg := newConfig()

	bot, err := telegram.NewBot(cfg.telegramToken)
	if err != nil {
		slog.Error("failed to create bot", slog.Any("error", err))
		return
	}

	vs := newVibeshare(cfg, bot.Sender())
	bot.HandleText(vs.TextHandler)
	bot.HandleCallback(vs.CallbackHandler)
	defer bot.Stop()

	slog.Info("Bot started")
	bot.Start()
}

func newConfig() *config {
	return &config{
		telegramToken:       os.Getenv("TELEGRAM_TOKEN"),
		spotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		youtubeAPIKey:       os.Getenv("YOUTUBE_API_KEY"),
	}
}

func newVibeshare(cfg *config, ts telegram.Sender) *vibeshare.Vibeshare {
	return vibeshare.NewVibeshare(&vibeshare.Input{
		MusicRegistry: music.NewRegistry(&music.RegistryInput{
			SpotifyClient: spotify.NewHTTPClient(&spotify.Credentials{
				ClientID:     cfg.spotifyClientID,
				ClientSecret: cfg.spotifyClientSecret,
			}),
			YandexClient: yandex.NewHTTPClient(),
			YoutubeClient: youtube.NewHTTPClient(
				cfg.youtubeAPIKey,
			),
		}),
		TelegramSender: ts,
	})
}
