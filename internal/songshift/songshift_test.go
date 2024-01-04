package songshift

import (
	"fmt"
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/ymusic"
	spotify_utils "github.com/GeorgeGorbanev/songshift/tests/utils/spotify"
	telegram_utils "github.com/GeorgeGorbanev/songshift/tests/utils/telegram"
	ymusic_utils "github.com/GeorgeGorbanev/songshift/tests/utils/ymusic"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestSongshift_HandleText(t *testing.T) {
	sampleSender := &telebot.User{Username: "sample_username"}

	tests := []struct {
		name             string
		inMsg            *telebot.Message
		expectedResponse *telegram.Message
	}{
		{
			name: "when no link given",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text:   "just message with no link",
			},
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "no track link found",
			},
		},
		{
			name: "when spotify track link given and track found",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text: fmt.Sprintf(
					"prfx https://open.spotify.com/track/%s?sample=query",
					spotify_utils.TrackFixtureMassiveAttackAngel.Track.ID,
				),
			},
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "https://music.yandex.com/album/35627/track/354093",
			},
		},
		{
			name: "when spotify track link given and track found but ymusic track not found",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text: fmt.Sprintf(
					"prfx https://open.spotify.com/track/%s?sample=query",
					spotify_utils.TrackFixtureMileyCyrusFlowers.Track.ID,
				),
			},
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "no ym track found",
			},
		},
		{
			name: "when spotify track link given and track not found",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text:   fmt.Sprintf("prfx https://open.spotify.com/track/%s?sample=query", "invalid_track_id"),
			},
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "track not found",
			},
		},
		{
			name: "when yandex music track link given and track found",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text: fmt.Sprintf(
					"prfx https://music.yandex.ru/album/3192570/track/%s?query=sample",
					ymusic_utils.TrackFixtureMassiveAttackAngel.ID,
				),
			},
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
			},
		},
		{
			name: "when yandex music track link given and track not found",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text:   fmt.Sprintf("prfx https://music.yandex.ru/album/3192570/track/%s", "0"),
			},
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "track not found in yandex music",
			},
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

			ss := NewSongshift(&Input{
				SpotifyClient:  spotifyClient,
				TelegramSender: senderMock,
				YmusicClient:   ymClient,
			})
			ss.HandleText(tt.inMsg)

			require.Equal(t, tt.expectedResponse, senderMock.Response)
		})
	}
}
