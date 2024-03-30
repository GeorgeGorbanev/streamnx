package convert_album

import (
	"os"
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/youtube"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/tucnak/telebot"
)

var (
	fixturesMap *utils.FixturesMap
	senderMock  *utils.TelegramSenderMock
	vs          *vibeshare.Vibeshare

	user = &telebot.User{
		Username: "sample_username",
	}
)

func TestMain(m *testing.M) {
	senderMock = utils.NewTelegramSenderMock()
	fixturesMap = &utils.FixturesMap{}

	spotifyAuthServerMock := utils.NewSpotifyAuthServerMock()
	spotifyAPIServerMock := utils.NewSpotifyAPIServerMock(fixturesMap)
	youtubeMockServer := utils.NewYoutubeAPIServerMock(fixturesMap)
	yandexMockServer := utils.NewYandexAPIServerMock(fixturesMap)

	spotifyClient := spotify.NewHTTPClient(
		&utils.SpotifyCredentials,
		spotify.WithAuthURL(spotifyAuthServerMock.URL),
		spotify.WithAPIURL(spotifyAPIServerMock.URL),
	)
	youtubeClient := youtube.NewHTTPClient(utils.YoutubeAPIKey, youtube.WithAPIURL(youtubeMockServer.URL))
	yandexClient := yandex.NewHTTPClient(yandex.WithAPIURL(yandexMockServer.URL))

	vs = vibeshare.NewVibeshare(&vibeshare.Input{
		MusicRegistry: music.NewRegistry(&music.RegistryInput{
			SpotifyClient: spotifyClient,
			YandexClient:  yandexClient,
			YoutubeClient: youtubeClient,
		}),
		TelegramSender: senderMock,
	})

	code := m.Run()

	spotifyAuthServerMock.Close()
	spotifyAPIServerMock.Close()
	yandexMockServer.Close()
	youtubeMockServer.Close()

	os.Exit(code)
}
