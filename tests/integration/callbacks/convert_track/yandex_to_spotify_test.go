package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/vibeshare/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackYandexToSpotify(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when yandex track link given and spotify track found",
			input: "cnvtr/ya/354093/sf",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://open.spotify.com/track/7uv632EkfwYhXoqf8rhYrg",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YandexTracks: map[string][]byte{
					"354093": fixture.Read("yandex/get_track_massive_attack_angel.json"),
				},
				SpotifySearchTracks: map[string][]byte{
					"artist:Massive Attack track:Angel": fixture.Read("spotify/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:  "when yandex track link given and spotify track not found",
			input: "cnvtr/ya/354093/sf",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Track not found in Spotify",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YandexTracks: map[string][]byte{
					"354093": fixture.Read("yandex/get_track_massive_attack_angel.json"),
				},
				SpotifySearchTracks: map[string][]byte{},
			},
		},
		{
			name:             "when yandex track not found",
			input:            "cnvtr/ya/354093/sf",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				YandexTracks:        map[string][]byte{},
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
