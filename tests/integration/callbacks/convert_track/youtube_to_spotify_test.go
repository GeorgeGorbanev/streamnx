package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackYoutubeToSpotify(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  utils.FixturesMap
	}{
		{
			name:         "when youtube track link given and spotify track found",
			input:        "cnvtr/yt/hbe3CQamF8k/sf",
			expectedText: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
			fixturesMap: utils.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"hbe3CQamF8k": fixture.Read("youtube/get_track_massive_attack_angel.json"),
				},
				SpotifySearchTracks: map[string][]byte{
					// TODO: fix query
					"artist:Massive Attack - Angel track:": fixture.Read("spotify/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when yotube track link given and spotify track not found",
			input:        "cnvtr/yt/hbe3CQamF8k/sf",
			expectedText: "Track not found in Spotify",
			fixturesMap: utils.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"hbe3CQamF8k": fixture.Read("youtube/get_track_massive_attack_angel.json"),
				},
				SpotifySearchTracks: map[string][]byte{},
			},
		},
		{
			name:         "when youtube track not found",
			input:        "cnvtr/yt/hbe3CQamF8k/sf",
			expectedText: "",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks:       map[string][]byte{},
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
