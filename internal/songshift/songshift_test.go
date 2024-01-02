package songshift

import (
	"fmt"
	"net/http/httptest"
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
		name                    string
		inMsg                   *telebot.Message
		yMusicMockServerCreator func(t *testing.T) *httptest.Server
		expectedResponse        *telegram.Message
	}{
		{
			name: "when no link given",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text:   "just message with no link",
			},
			yMusicMockServerCreator: ymusic_utils.NewAPISearchServerMock,
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "no track link found",
			},
		},
		{
			name: "when spotify track link given and track found",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text:   fmt.Sprintf("prfx https://open.spotify.com/track/%s?sample=query", spotify_utils.SampleTrack.ID),
			},
			yMusicMockServerCreator: ymusic_utils.NewAPISearchServerMock,
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "https://music.yandex.com/album/35627/track/354093",
			},
		},
		{
			name: "when spotify track link given and track found but ymusic track not found",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text:   fmt.Sprintf("prfx https://open.spotify.com/track/%s?sample=query", spotify_utils.SampleYMusicNotFoundTrack.ID),
			},
			yMusicMockServerCreator: ymusic_utils.NewAPISearchServerMock,
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
			yMusicMockServerCreator: ymusic_utils.NewAPISearchServerMock,
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "track not found",
			},
		},
		{
			name: "when yandex music track link given and track found",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text:   fmt.Sprintf("prfx https://music.yandex.ru/album/3192570/track/%s", ymusic_utils.SampleTrackID),
			},
			yMusicMockServerCreator: ymusic_utils.NewAPIGetTrackServerMock,
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "The Raveonettes - Come On Santa",
			},
		},
		{
			name: "when yandex music track link given and track not found",
			inMsg: &telebot.Message{
				Sender: sampleSender,
				Text:   fmt.Sprintf("prfx https://music.yandex.ru/album/3192570/track/%s", "0"),
			},
			yMusicMockServerCreator: ymusic_utils.NewAPIGetTrackServerMock,
			expectedResponse: &telegram.Message{
				To:   sampleSender,
				Text: "track not found",
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

			require.Equal(t, ymusic_utils.SampleSearchQuery, spotify_utils.SampleTrack.Title())

			ymusicMockServer := tt.yMusicMockServerCreator(t)
			defer ymusicMockServer.Close()

			ymClient := ymusic.NewClient(ymusic.WithAPIURL(ymusicMockServer.URL))

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
