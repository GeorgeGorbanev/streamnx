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

func TestMessage_SpotifyTrack(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedResponse string
	}{
		{
			name: "when spotify track link given and track found",
			input: fmt.Sprintf(
				"prfx https://open.spotify.com/track/%s?sample=query",
				spotify_utils.TrackFixtureMassiveAttackAngel.Track.ID,
			),
			expectedResponse: "https://music.yandex.com/album/35627/track/354093",
		},
		{
			name:             "when spotify track link given and track not found",
			input:            "https://open.spotify.com/track/invalid_track_id?sample=query",
			expectedResponse: "track not found",
		},
		{
			name: "when spotify track link given, track found and ymusic track found, but artist name not match",
			input: fmt.Sprintf(
				"prfx https://open.spotify.com/track/%s?sample=query",
				spotify_utils.TrackFixtureMileyCyrusFlowers.Track.ID,
			),
			expectedResponse: "no ym track found",
		},
		{
			name: "when spotify track link given, ymusic track found and artist name not match, but match in translit",
			input: fmt.Sprintf(
				"prfx https://open.spotify.com/track/%s?sample=query",
				spotify_utils.TrackFixtureZemfiraIskala.Track.ID,
			),
			expectedResponse: fmt.Sprintf(
				"https://music.yandex.com/album/%d/track/%s",
				ymusic_utils.TrackFixtureZemfiraIskala.Track.Albums[0].ID,
				ymusic_utils.TrackFixtureZemfiraIskala.Track.IDString(),
			),
		},
		{
			name: "when spotify track link given, track found, ymusic track not found, but found in translit",
			input: fmt.Sprintf(
				"prfx https://open.spotify.com/track/%s?sample=query",
				spotify_utils.TrackFixtureNadezhdaKadyshevaShirokaReka.Track.ID,
			),
			expectedResponse: fmt.Sprintf(
				"https://music.yandex.com/album/%d/track/%s",
				ymusic_utils.TrackFixtureNadezhdaKadyshevaShirokaReka.Track.Albums[0].ID,
				ymusic_utils.TrackFixtureNadezhdaKadyshevaShirokaReka.Track.IDString(),
			),
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
