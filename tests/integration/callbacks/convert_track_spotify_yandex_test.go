package callbacks

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/yandex"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackSpotifyToYandex(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  utils.FixturesMap
	}{
		{
			name:         "when spotify track link given and yandex track found",
			input:        "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ya",
			expectedText: "https://music.yandex.com/album/35627/track/354093",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"massive attack – angel": fixture.Read("yandex/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when spotify track link given and yandex track not found",
			input:        "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ya",
			expectedText: "Track not found in Yandex",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when yandex track not found",
			input:        "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ya",
			expectedText: "",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks:      map[string][]byte{},
				YandexSearchTracks: map[string][]byte{},
			},
		},
		{
			name:         "when spotify track link given, track found and yandex track found, but artist name not match",
			input:        "cnvtr/sf/7DSAEUvxU8FajXtRloy8M0/ya",
			expectedText: "Track not found in Yandex",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7DSAEUvxU8FajXtRloy8M0": fixture.Read("spotify/get_track_miley_cyrus_flowers.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"massive attack – angel": fixture.Read("yandex/search_track_miley_cyrus_flowers.json"),
				},
			},
		},
		{
			name:         "when spotify track link given, yandex track found and artist name not match, but match in translit",
			input:        "cnvtr/sf/3NHSz1GyC5IeK1soZSjIIX/ya",
			expectedText: "https://music.yandex.com/album/81431/track/732401",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"3NHSz1GyC5IeK1soZSjIIX": fixture.Read("spotify/get_track_zemfira_iskala.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"zemfira – искала": fixture.Read("yandex/search_track_zemfira_iskala.json"),
				},
			},
		},
		{
			name:         "when spotify track link given, track found, yandex track not found, but found in translit",
			input:        "cnvtr/sf/2sP5VgY8PWb6c9DhgZEpSv/ya",
			expectedText: "https://music.yandex.com/album/4058886/track/33223088",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"2sP5VgY8PWb6c9DhgZEpSv": fixture.Read("spotify/get_track_nadezhda_kadysheva_shiroka_reka.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"надежда кадышева – широка река": fixture.Read("yandex/search_track_nadezhda_kadysheva_shiroka_reka.json"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spotifyAuthServerMock := utils.NewSpotifyAuthServerMock(t)
			defer spotifyAuthServerMock.Close()

			spotifyAPIServerMock := utils.NewSpotifyAPIServerMock(t, tt.fixturesMap)
			defer spotifyAPIServerMock.Close()

			senderMock := utils.NewTelegramSenderMock()
			spotifyClient := spotify.NewHTTPClient(
				&utils.SpotifyCredentials,
				spotify.WithAuthURL(spotifyAuthServerMock.URL),
				spotify.WithAPIURL(spotifyAPIServerMock.URL),
			)

			yandexMockServer := utils.NewYandexAPIServerMock(t, tt.fixturesMap)
			defer yandexMockServer.Close()

			yandexClient := yandex.NewHTTPClient(yandex.WithAPIURL(yandexMockServer.URL))

			vs := vibeshare.NewVibeshare(&vibeshare.Input{
				MusicRegistry: music.NewRegistry(&music.RegistryInput{
					SpotifyClient: spotifyClient,
					YandexClient:  yandexClient,
				}),
				TelegramSender: senderMock,
			})

			user := &telebot.User{
				Username: "sample_username",
			}

			callback := telebot.Callback{
				Sender: user,
				Data:   tt.input,
			}

			vs.CallbackHandler(&callback)

			if tt.expectedText == "" {
				require.Nil(t, senderMock.Response)
			} else {
				require.NotNil(t, senderMock.Response)
				require.Equal(t, user, senderMock.Response.To)
				require.Equal(t, tt.expectedText, senderMock.Response.Text)
			}
		})
	}
}
