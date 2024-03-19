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

func TestMessage_YandexAlbum(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedResponse string
	}{
		{
			name: "when yandex music album link given and album found",
			input: fmt.Sprintf(
				"https://music.yandex.com/album/%s",
				ymusic_utils.AlbumFixtureRadioheadAmnesiac.ID,
			),
			expectedResponse: "https://open.spotify.com/album/1HrMmB5useeZ0F5lHrMvl0",
		},
		{
			name: "when yandex music album link given and yandex music album not found",
			input: fmt.Sprintf(
				"https://music.yandex.com/album/%s",
				"0",
			),
			expectedResponse: "no yandex music album found",
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
