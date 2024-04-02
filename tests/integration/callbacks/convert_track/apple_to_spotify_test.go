package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackAppleToSpotify(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  fixture.FixturesMap
	}{
		{
			name:         "when apple track link given and spotify track found",
			input:        "cnvtr/ap/724466660/sf",
			expectedText: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
			fixturesMap: fixture.FixturesMap{
				AppleTracks: map[string][]byte{
					"724466660": fixture.Read("apple/get_track_massive_attack_angel.json"),
				},
				SpotifySearchTracks: map[string][]byte{
					"artist:Massive Attack track:Angel": fixture.Read("spotify/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when apple track link given and spotify track not found",
			input:        "cnvtr/ap/724466660/sf",
			expectedText: "Track not found in Spotify",
			fixturesMap: fixture.FixturesMap{
				AppleTracks: map[string][]byte{
					"724466660": fixture.Read("apple/get_track_massive_attack_angel.json"),
				},
				SpotifySearchTracks: map[string][]byte{},
			},
		},
		{
			name:         "when apple track not found",
			input:        "cnvtr/ap/724466660/sf",
			expectedText: "",
			fixturesMap: fixture.FixturesMap{
				AppleTracks:         map[string][]byte{},
				SpotifySearchTracks: map[string][]byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixturesMap.Merge(&tt.fixturesMap)
			defer fixturesMap.Reset()
			defer senderMock.Reset()

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
