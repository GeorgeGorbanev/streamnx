package songshift

import (
	"fmt"
	"testing"

	"github.com/GeorgeGorbanev/songshift/internal/songshift/spotify"
	"github.com/GeorgeGorbanev/songshift/internal/songshift/telegram"
	spotify_utils "github.com/GeorgeGorbanev/songshift/tests/utils/spotify"
	telegram_utils "github.com/GeorgeGorbanev/songshift/tests/utils/telegram"
	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestSongshift_HandleText(t *testing.T) {
	tests := []struct {
		name             string
		inMsg            *telebot.Message
		expectedResponse *telegram.Message
	}{
		{
			name: "when no link given",
			inMsg: &telebot.Message{
				Sender: &telebot.User{
					Username: "sample_username",
				},
				Text: "just message with no link",
			},
			expectedResponse: &telegram.Message{
				To:   &telebot.User{Username: "sample_username"},
				Text: "no track link found",
			},
		},
		{
			name: "when spotify track link given and track found",
			inMsg: &telebot.Message{
				Sender: &telebot.User{
					Username: "sample_username",
				},
				Text: fmt.Sprintf("prfx https://open.spotify.com/track/%s?sample=query", spotify_utils.SampleTrack.ID),
			},
			expectedResponse: &telegram.Message{
				To:   &telebot.User{Username: "sample_username"},
				Text: `Track: "Artist One – Track One"`,
			},
		},
		{
			name: "when spotify track link given and track not found",
			inMsg: &telebot.Message{
				Sender: &telebot.User{
					Username: "sample_username",
				},
				Text: fmt.Sprintf("prfx https://open.spotify.com/track/%s?sample=query", "invalid_track_id"),
			},
			expectedResponse: &telegram.Message{
				To:   &telebot.User{Username: "sample_username"},
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
			ss := NewSongshift(spotifyClient, senderMock)
			ss.HandleText(tt.inMsg)

			require.Equal(t, senderMock.Response, tt.expectedResponse)
		})
	}
}
