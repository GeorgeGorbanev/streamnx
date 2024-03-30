package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackSpotifyToYoutube(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  utils.FixturesMap
	}{
		{
			name:         "when spotify track link given and youtube track found",
			input:        "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/yt",
			expectedText: "https://www.youtube.com/watch?v=hbe3CQamF8k",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
				YoutubeSearchTracks: map[string][]byte{
					"Massive Attack – Angel": fixture.Read("youtube/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when spotify track link given and youtube track not found",
			input:        "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/yt",
			expectedText: "Track not found in Youtube",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
				YoutubeSearchTracks: map[string][]byte{},
			},
		},
		{
			name:         "when yandex track not found",
			input:        "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/yt",
			expectedText: "",
			fixturesMap: utils.FixturesMap{
				SpotifyTracks:       map[string][]byte{},
				YoutubeSearchTracks: map[string][]byte{},
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
