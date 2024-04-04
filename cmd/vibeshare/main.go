package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/apple"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/youtube"

	"github.com/joho/godotenv"
)

type config struct {
	telegramToken       string
	spotifyClientID     string
	spotifyClientSecret string
	youtubeAPIKey       string
	feedbackToken       string
	feedbackReceiverID  int
}

func main() {
	setupLogs()
	cfg, err := loadConfig()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return
	}

	vs, err := newVibeshare(cfg)
	if err != nil {
		slog.Error("failed to create vibeshare", slog.Any("error", err))
		return
	}
	defer vs.Stop()

	vs.Run()
}

func setupLogs() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{AddSource: true},
			),
		),
	)
}

func loadConfig() (*config, error) {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, err
		}
	}

	fbReceiverID, err := strconv.Atoi(os.Getenv("FEEDBACK_RECEIVER_ID"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse feedback receiver: %w", err)
	}

	return &config{
		telegramToken:       os.Getenv("TELEGRAM_TOKEN"),
		spotifyClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifyClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		youtubeAPIKey:       os.Getenv("YOUTUBE_API_KEY"),
		feedbackToken:       os.Getenv("FEEDBACK_TOKEN"),
		feedbackReceiverID:  fbReceiverID,
	}, nil
}

func newVibeshare(cfg *config) (vibeshare.Vibeshare, error) {
	return vibeshare.NewVibeshare(&vibeshare.Input{
		VibeshareBotToken:  cfg.telegramToken,
		FeedbackBotToken:   cfg.feedbackToken,
		FeedbackReceiverID: cfg.feedbackReceiverID,
		MusicRegistry: music.NewRegistry(&music.RegistryInput{
			AppleClient: apple.NewHTTPClient(),
			SpotifyClient: spotify.NewHTTPClient(&spotify.Credentials{
				ClientID:     cfg.spotifyClientID,
				ClientSecret: cfg.spotifyClientSecret,
			}),
			YandexClient: yandex.NewHTTPClient(),
			YoutubeClient: youtube.NewHTTPClient(
				cfg.youtubeAPIKey,
			),
		}),
	})
}
