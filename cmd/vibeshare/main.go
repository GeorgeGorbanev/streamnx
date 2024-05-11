package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/GeorgeGorbanev/vibeshare/internal/apple"
	"github.com/GeorgeGorbanev/vibeshare/internal/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/translator"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/youtube"

	"github.com/joho/godotenv"
)

type config struct {
	googleCloudProjectID       string
	googleTranslatorAPIKeyJSON string
	telegramToken              string
	spotifyClientID            string
	spotifyClientSecret        string
	youtubeAPIKey              string
	feedbackToken              string
	feedbackReceiverID         int
}

func main() {
	setupLogs()
	cfg, err := loadConfig()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		return
	}

	ctx := context.Background()

	vs, err := newVibeshare(ctx, cfg)
	if err != nil {
		slog.Error("failed to start vibeshare", slog.Any("error", err))
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
		telegramToken:              os.Getenv("TELEGRAM_TOKEN"),
		googleCloudProjectID:       os.Getenv("GOOGLE_CLOUD_PROJECT_ID"),
		googleTranslatorAPIKeyJSON: os.Getenv("GOOGLE_TRANSLATOR_KEY_JSON"),
		spotifyClientID:            os.Getenv("SPOTIFY_CLIENT_ID"),
		spotifyClientSecret:        os.Getenv("SPOTIFY_CLIENT_SECRET"),
		youtubeAPIKey:              os.Getenv("YOUTUBE_API_KEY"),
		feedbackToken:              os.Getenv("FEEDBACK_TOKEN"),
		feedbackReceiverID:         fbReceiverID,
	}, nil
}

func newVibeshare(ctx context.Context, cfg *config) (*vibeshare.Vibeshare, error) {
	googleClient, err := translator.NewGoogleClient(ctx, &translator.GoogleCredentials{
		APIKeyJSON: cfg.googleTranslatorAPIKeyJSON,
		ProjectID:  cfg.googleCloudProjectID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create google translator client: %w", err)
	}

	return vibeshare.NewVibeshare(&vibeshare.Input{
		VibeshareBotToken:  cfg.telegramToken,
		FeedbackBotToken:   cfg.feedbackToken,
		FeedbackReceiverID: cfg.feedbackReceiverID,
		MusicRegistry: music.NewRegistry(&music.RegistryInput{
			AppleClient:  apple.NewHTTPClient(),
			YandexClient: yandex.NewHTTPClient(),
			YoutubeClient: youtube.NewHTTPClient(
				cfg.youtubeAPIKey,
			),
			SpotifyClient: spotify.NewHTTPClient(&spotify.Credentials{
				ClientID:     cfg.spotifyClientID,
				ClientSecret: cfg.spotifyClientSecret,
			}),
			Translator: googleClient,
		}),
	})
}
