package convert_track

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertTrackAppleToYandex(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  fixture.FixturesMap
	}{
		{
			name:         "when apple track link given and yandex track found",
			input:        "cnvtr/ap/724466660/ya",
			expectedText: "https://music.yandex.com/album/35627/track/354093",
			fixturesMap: fixture.FixturesMap{
				AppleTracks: map[string][]byte{
					"724466660": fixture.Read("apple/get_track_massive_attack_angel.json"),
				},
				YandexSearchTracks: map[string][]byte{
					"massive attack – angel": fixture.Read("yandex/search_track_massive_attack_angel.json"),
				},
			},
		},
		{
			name:         "when apple track link given and yandex track not found",
			input:        "cnvtr/ap/724466660/ya",
			expectedText: "Track not found in Yandex",
			fixturesMap: fixture.FixturesMap{
				AppleTracks: map[string][]byte{
					"724466660": fixture.Read("apple/get_track_massive_attack_angel.json"),
				},
				YandexSearchTracks: map[string][]byte{},
			},
		},
		{
			name:         "when apple track not found",
			input:        "cnvtr/ap/724466660/ya",
			expectedText: "",
			fixturesMap: fixture.FixturesMap{
				AppleTracks:        map[string][]byte{},
				YandexSearchTracks: map[string][]byte{},
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
