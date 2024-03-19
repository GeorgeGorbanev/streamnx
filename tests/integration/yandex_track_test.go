package integration

import (
	"fmt"
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/spotify"
	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/ymusic"
	spotify_utils "github.com/GeorgeGorbanev/vibeshare/tests/utils/spotify"
	telegram_utils "github.com/GeorgeGorbanev/vibeshare/tests/utils/telegram"
	ymusic_utils "github.com/GeorgeGorbanev/vibeshare/tests/utils/ymusic"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestMessage_YandexTrack(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedResponse string
	}{
		{
			name: "when yandex music track link given and track found with .ru",
			input: fmt.Sprintf(
				"prfx https://music.yandex.ru/album/3192570/track/%s?query=sample",
				ymusic_utils.TrackFixtureMassiveAttackAngel.ID,
			),
			expectedResponse: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name: "when yandex music track link given and track found with .com",
			input: fmt.Sprintf(
				"https://music.yandex.com/album/3192570/track/%s",
				ymusic_utils.TrackFixtureMassiveAttackAngel.ID,
			),
			expectedResponse: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
		},
		{
			name:             "when yandex music track link given and track not found",
			input:            fmt.Sprintf("prfx https://music.yandex.ru/album/3192570/track/%s", "0"),
			expectedResponse: "track not found in yandex music",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAuthServer := spotify_utils.NewAuthServerMock(t)
			defer mockAuthServer.Close()

			mockAPIServer := spotify_utils.NewAPIServerMock(t)
			defer mockAPIServer.Close()

			senderMock := telegram_utils.NewTelegramSenderMock()
			spotifyClient := spotify.NewClient(
				&spotify_utils.SampleCredentials,
				spotify.WithAuthURL(mockAuthServer.URL),
				spotify.WithAPIURL(mockAPIServer.URL),
			)

			yMusicMockServer := ymusic_utils.NewAPIServerMock(t)
			defer yMusicMockServer.Close()

			ymClient := ymusic.NewClient(ymusic.WithAPIURL(yMusicMockServer.URL))

			vs := vibeshare.NewVibeshare(&vibeshare.Input{
				SpotifyClient:  spotifyClient,
				TelegramSender: senderMock,
				YmusicClient:   ymClient,
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
