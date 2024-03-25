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

func TestMessage_SpotifyTrack(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedResponse string
		fixturesMap      fixturesMap
	}{
		{
			name:             "when spotify track link given and track found",
			input:            "prfx https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg?sample=query",
			expectedResponse: "https://music.yandex.com/album/35627/track/354093",
			fixturesMap: fixturesMap{
				spotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
				yandexSearchTracks: map[string][]byte{
					"massive attack – angel": fixture.Read("yandex/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:             "when spotify track link given and track not found",
			input:            "https://open.spotify.com/track/invalid_track_id?sample=query",
			expectedResponse: "failed to convert",
			fixturesMap:      fixturesMap{},
		},
		{
			name:             "when spotify track link given, track found and yandex track found, but artist name not match",
			input:            "prfx https://open.spotify.com/track/7DSAEUvxU8FajXtRloy8M0?sample=query",
			expectedResponse: "failed to convert",
			fixturesMap: fixturesMap{
				spotifyTracks: map[string][]byte{
					"7DSAEUvxU8FajXtRloy8M0": fixture.Read("spotify/get_track_miley_cyrus_flowers.json"),
				},
				yandexSearchTracks: map[string][]byte{
					"massive attack – angel": fixture.Read("yandex/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:             "when spotify track link given, yandex track found and artist name not match, but match in translit",
			input:            "prfx https://open.spotify.com/track/3NHSz1GyC5IeK1soZSjIIX?sample=query",
			expectedResponse: "https://music.yandex.com/album/81431/track/732401",
			fixturesMap: fixturesMap{
				spotifyTracks: map[string][]byte{
					"3NHSz1GyC5IeK1soZSjIIX": fixture.Read("spotify/get_track_zemfira_iskala.json"),
				},
				yandexSearchTracks: map[string][]byte{
					"zemfira – искала": fixture.Read("yandex/search_track_zemfira_iskala.json"),
				},
			},
		},
		{
			name:             "when spotify track link given, track found, yandex track not found, but found in translit",
			input:            "prfx https://open.spotify.com/track/2sP5VgY8PWb6c9DhgZEpSv?sample=query",
			expectedResponse: "https://music.yandex.com/album/4058886/track/33223088",
			fixturesMap: fixturesMap{
				spotifyTracks: map[string][]byte{
					"2sP5VgY8PWb6c9DhgZEpSv": fixture.Read("spotify/get_track_nadezhda_kadysheva_shiroka_reka.json"),
				},
				yandexSearchTracks: map[string][]byte{
					"надежда кадышева – широка река": fixture.Read("yandex/search_track_nadezhda_kadysheva_shiroka_reka.json"),
				},
			},
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

			vs.TextHandler(msg)

			require.Equal(t, user, senderMock.Response.To)
			require.Equal(t, tt.expectedResponse, senderMock.Response.Text)
		})
	}
}
