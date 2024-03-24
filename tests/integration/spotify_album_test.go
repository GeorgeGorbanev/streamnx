package integration

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/converter"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	telegram_utils "github.com/GeorgeGorbanev/vibeshare/tests/utils/telegram"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestMessage_SpotifyAlbum(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedResponse string
		fixturesMap      fixturesMap
	}{
		{
			name:             "when spotify album link given and album found",
			input:            "https://open.spotify.com/album/1HrMmB5useeZ0F5lHrMvl0",
			expectedResponse: "https://music.yandex.com/album/3389008",
			fixturesMap: fixturesMap{
				spotifyAlbums: map[string][]byte{
					"1HrMmB5useeZ0F5lHrMvl0": fixture.Read("spotify/get_album_radiohead_amnesiac.json"),
				},
				yandexSearchAlbums: map[string][]byte{
					"radiohead – amnesiac": fixture.Read("yandex/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:             "when spotify album link given and spotify album not found",
			input:            "https://open.spotify.com/album/0",
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
			spotifyClient := spotify.NewHTTPClient(
				&sampleCredentials,
				spotify.WithAuthURL(spotifyAuthServerMock.URL),
				spotify.WithAPIURL(spotifyAPIServerMock.URL),
			)

			yandexMockServer := newYandexAPIServerMock(t, tt.fixturesMap)
			defer yandexMockServer.Close()

			yandexClient := yandex.NewHTTPClient(yandex.WithAPIURL(yandexMockServer.URL))

			c := converter.NewConverter(&converter.Input{
				SpotifyClient: spotifyClient,
				YandexClient:  yandexClient,
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
