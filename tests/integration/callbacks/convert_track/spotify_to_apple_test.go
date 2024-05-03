package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackSpotifyToApple(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when spotify track link given and apple track found",
			input: "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ap",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.apple.com/us/album/angel/724466069?i=724466660",
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
				AppleSearchTracks: map[string][]byte{
					"Massive Attack Angel": fixture.Read("apple/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:  "when spotify track link given and apple track not found",
			input: "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ap",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Track not found in Apple",
				},
			},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks: map[string][]byte{
					"7uv632EkfwYhXoqf8rhYrg": fixture.Read("spotify/get_track_massive_attack_angel.json"),
				},
				AppleSearchTracks: map[string][]byte{},
			},
		},
		{
			name:             "when spotify track not found",
			input:            "cnvtr/sf/7uv632EkfwYhXoqf8rhYrg/ap",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				SpotifyTracks:     map[string][]byte{},
				AppleSearchTracks: map[string][]byte{},
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
