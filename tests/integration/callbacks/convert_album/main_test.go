package convert_album

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/streaminx"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/tucnak/telebot"
)

var (
	fixturesMap    *fixture.FixturesMap
	senderMock     *utils.TelegramSenderMock
	translatorMock utils.TranslatorMock
	vs             *vibeshare.Vibeshare

	user = &telebot.User{
		Username: "sample_username",
	}
)

func TestMain(m *testing.M) {
	senderMock = utils.NewTelegramSenderMock()
	fixturesMap = &fixture.FixturesMap{}

	appleAPIServerMock := utils.NewAppleAPIServerMock(fixturesMap)
	appleWebPlayerServerMock := utils.NewAppleWebPlayerServerMock()
	spotifyAuthServerMock := utils.NewSpotifyAuthServerMock()
	spotifyAPIServerMock := utils.NewSpotifyAPIServerMock(fixturesMap)
	youtubeMockServer := utils.NewYoutubeAPIServerMock(fixturesMap)
	yandexMockServer := utils.NewYandexAPIServerMock(fixturesMap)

	streaminxRegistry, err := streaminx.NewRegistry(
		context.Background(),
		streaminx.Credentials{
			SpotifyClientID:     utils.SpotifyCredentials.ClientID,
			SpotifyClientSecret: utils.SpotifyCredentials.ClientSecret,
			YoutubeAPIKey:       utils.YoutubeAPIKey,
		},
		streaminx.WithAppleAPIURL(appleAPIServerMock.URL),
		streaminx.WithAppleWebPlayerURL(appleWebPlayerServerMock.URL),
		streaminx.WithSpotifyAuthURL(spotifyAuthServerMock.URL),
		streaminx.WithSpotifyAPIURL(spotifyAPIServerMock.URL),
		streaminx.WithYandexAPIURL(yandexMockServer.URL),
		streaminx.WithYoutubeAPIURL(youtubeMockServer.URL),
		streaminx.WithTranslator(&translatorMock),
	)
	if err != nil {
		slog.Error("failed to build streaminx registry", slog.Any("error", err))
		os.Exit(1)
	}

	app, err := vibeshare.NewVibeshare(
		&vibeshare.Input{
			StreaminxRegistry: streaminxRegistry,
		},
		vibeshare.WithVibeshareSender(senderMock),
	)
	if err != nil {
		slog.Error("failed to build vibeshare", slog.Any("error", err))
		os.Exit(1)
	}
	vs = app

	code := m.Run()

	appleAPIServerMock.Close()
	appleWebPlayerServerMock.Close()
	spotifyAuthServerMock.Close()
	spotifyAPIServerMock.Close()
	yandexMockServer.Close()
	youtubeMockServer.Close()

	os.Exit(code)
}
