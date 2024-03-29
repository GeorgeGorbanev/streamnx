package texts

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/music"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestText_SpotifyTrackLink(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedText        string
		expectedReplyMarkup *telebot.ReplyMarkup
		fixturesMap         utils.FixturesMap
	}{
		{
			name:         "when spotify track link given and track found",
			input:        "prfx https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg?sample=query",
			expectedText: "Select target link provider",
			expectedReplyMarkup: &telebot.ReplyMarkup{
				InlineKeyboard: [][]telebot.InlineButton{
					{
						{
							Text: "Yandex",
							Data: "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ya",
						},
						{
							Text: "Youtube",
							Data: "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/yt",
						},
					},
				},
			},
			fixturesMap: utils.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when spotify track link given and track not found",
			input:        "https://open.spotify.com/track/0?sample=query",
			expectedText: "Link is invalid",
			fixturesMap:  utils.FixturesMap{},
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

			vs := vibeshare.NewVibeshare(&vibeshare.Input{
				MusicRegistry: music.NewRegistry(&music.RegistryInput{
					SpotifyClient: spotifyClient,
				}),
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
			require.Equal(t, tt.expectedText, senderMock.Response.Text)
			require.Equal(t, tt.expectedReplyMarkup, senderMock.Response.ReplyMarkup)
		})
	}
}
