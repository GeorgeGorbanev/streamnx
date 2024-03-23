package integration

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/converter"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/ymusic"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	telegram_utils "github.com/GeorgeGorbanev/vibeshare/tests/utils/telegram"
	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestMessage_YandexAlbum(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedResponse string
		fixturesMap      fixturesMap
	}{
		{
			name:             "when yandex music album link given and album found",
			input:            "https://music.yandex.com/album/3389008",
			expectedResponse: "https://open.spotify.com/album/1HrMmB5useeZ0F5lHrMvl0",
			fixturesMap: fixturesMap{
				yandexAlbums: map[string][]byte{
					"3389008": fixture.Read("yandex/get_album_radiohead_amnesiac.json"),
				},
				spotifySearchAlbums: map[string][]byte{
					"artist:Radiohead album:Amnesiac": fixture.Read("spotify/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:             "when yandex music album link given and yandex music album not found",
			input:            "https://music.yandex.com/album/0",
			expectedResponse: "failed to convert",
			fixturesMap:      fixturesMap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spotifyAuthServerMock := newSpotifyAuthServerMock(t)
			defer spotifyAuthServerMock.Close()

			spotifyAPIServerMock := newSpotifyAPIServerMock(t, tt.fixturesMap)
			defer spotifyAPIServerMock.Close()

			senderMock := telegram_utils.NewTelegramSenderMock()
			spotifyClient := spotify.NewClient(
				&sampleCredentials,
				spotify.WithAuthURL(spotifyAuthServerMock.URL),
				spotify.WithAPIURL(spotifyAPIServerMock.URL),
			)

			yMusicMockServer := newYandexAPIServerMock(t, tt.fixturesMap)
			defer yMusicMockServer.Close()

			ymClient := ymusic.NewClient(ymusic.WithAPIURL(yMusicMockServer.URL))

			c := converter.NewConverter(&converter.Input{
				SpotifyClient: spotifyClient,
				YandexClient:  ymClient,
			})
			vs := vibeshare.NewVibeshare(&vibeshare.Input{
				Converter:      c,
				TelegramSender: senderMock,
			})

			user := &telebot.User{
				Username: "sample_username",
			}
			msg := &telebot.Message{
				Sender: user,
				Text:   tt.input,
			}

			vs.HandleText(msg)

			require.Equal(t, user, senderMock.Response.To)
			require.Equal(t, tt.expectedResponse, senderMock.Response.Text)
		})
	}
}
