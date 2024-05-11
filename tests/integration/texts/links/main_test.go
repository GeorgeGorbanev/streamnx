package links

import (
	"os"
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/apple"
	"github.com/GeorgeGorbanev/vibeshare/internal/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/youtube"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/tucnak/telebot"
)

var (
	fixturesMap *fixture.FixturesMap
	senderMock  *utils.TelegramSenderMock
	vs          *vibeshare.Vibeshare

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

	appleClient := apple.NewHTTPClient(
		apple.WithAPIURL(appleAPIServerMock.URL),
		apple.WithWebPlayerURL(appleWebPlayerServerMock.URL),
	)
	spotifyClient := spotify.NewHTTPClient(
		&utils.SpotifyCredentials,
		spotify.WithAuthURL(spotifyAuthServerMock.URL),
		spotify.WithAPIURL(spotifyAPIServerMock.URL),
	)
	youtubeClient := youtube.NewHTTPClient(utils.YoutubeAPIKey, youtube.WithAPIURL(youtubeMockServer.URL))
	yandexClient := yandex.NewHTTPClient(yandex.WithAPIURL(yandexMockServer.URL))

	app, err := vibeshare.NewVibeshare(&vibeshare.Input{
		MusicRegistry: music.NewRegistry(&music.RegistryInput{
			AppleClient:   appleClient,
			SpotifyClient: spotifyClient,
			YandexClient:  yandexClient,
			YoutubeClient: youtubeClient,
		}),
	}, vibeshare.WithVibeshareSender(senderMock))
	if err != nil {
		panic(err)
	}
	vs = app

	code := m.Run()

	spotifyAuthServerMock.Close()
	spotifyAPIServerMock.Close()
	yandexMockServer.Close()
	youtubeMockServer.Close()

	os.Exit(code)
}
