package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackYoutubeToSpotify(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedText     string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:         "when youtube track link given and spotify track found",
			input:        "cnvtr/yt/hbe3CQamF8k/sf",
			expectedText: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"hbe3CQamF8k": fixture.Read("youtube/get_track_massive_attack_angel.json"),
				},
				SpotifySearchTracks: map[string][]byte{
					"artist:Massive Attack track:Angel": fixture.Read("spotify/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:  "when youtube track link given and spotify track not found",
			input: "cnvtr/yt/hbe3CQamF8k/sf",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Track not found in Spotify",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks: map[string][]byte{
					"hbe3CQamF8k": fixture.Read("youtube/get_track_massive_attack_angel.json"),
				},
				SpotifySearchTracks: map[string][]byte{},
			},
		},
		{
			name:             "when youtube track not found",
			input:            "cnvtr/yt/hbe3CQamF8k/sf",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				YoutubeTracks:       map[string][]byte{},
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

			require.Equal(t, tt.expectedMessages, senderMock.AllSent)
		})
	}
}
