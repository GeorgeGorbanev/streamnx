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

func TestMessage_YandexTrack(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedResponse string
		fixturesMap      fixturesMap
	}{
		{
			name:             "when yandex music track link given and track found with .ru",
			input:            "prfx https://music.yandex.ru/album/3192570/track/354093?query=sample",
			expectedResponse: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
			fixturesMap: fixturesMap{
				yandexTracks: map[string][]byte{
					"354093": fixture.Read("yandex/get_track_massive_attack_angel.json"),
				},
				spotifySearchTracks: map[string][]byte{
					"artist:Massive Attack track:Angel": fixture.Read("spotify/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:             "when yandex music track link given and track found with .com",
			input:            "https://music.yandex.com/album/3192570/track/354093",
			expectedResponse: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
			fixturesMap: fixturesMap{
				yandexTracks: map[string][]byte{
					"354093": fixture.Read("yandex/get_track_massive_attack_angel.json"),
				},
				spotifySearchTracks: map[string][]byte{
					"artist:Massive Attack track:Angel": fixture.Read("spotify/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:             "when yandex music track link given and track not found",
			input:            "prfx https://music.yandex.ru/album/3192570/track/0",
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

			yandexMockServer := newYandexAPIServerMock(t, tt.fixturesMap)
			defer yandexMockServer.Close()

			ymClient := yandex.NewClient(yandex.WithAPIURL(yandexMockServer.URL))

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
