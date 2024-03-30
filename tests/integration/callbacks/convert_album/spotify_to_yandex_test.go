package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"
	"github.com/GeorgeGorbanev/vibeshare/tests/utils"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumSpotifyToYandex(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  utils.FixturesMap
	}{
		{
			name:         "when spotify album link given and yandex album found",
			input:        "cnval/sf/1HrMmB5useeZ0F5lHrMvl0/ya",
			expectedText: "https://music.yandex.com/album/3389008",
			fixturesMap: utils.FixturesMap{
				SpotifyAlbums: map[string][]byte{
					"1HrMmB5useeZ0F5lHrMvl0": fixture.Read("spotify/get_album_radiohead_amnesiac.json"),
				},
				YandexSearchAlbums: map[string][]byte{
					"radiohead – amnesiac": fixture.Read("yandex/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:         "when spotify album link given and yandex album not found",
			input:        "cnval/sf/1HrMmB5useeZ0F5lHrMvl0/ya",
			expectedText: "Album not found in Yandex",
			fixturesMap: utils.FixturesMap{
				SpotifyAlbums: map[string][]byte{
					"1HrMmB5useeZ0F5lHrMvl0": fixture.Read("spotify/get_album_radiohead_amnesiac.json"),
				},
				YandexSearchAlbums: map[string][]byte{},
			},
		},
		{
			name:         "when spotify album not found",
			input:        "cnval/sf/1HrMmB5useeZ0F5lHrMvl0/ya",
			expectedText: "",
			fixturesMap: utils.FixturesMap{
				SpotifyAlbums:      map[string][]byte{},
				YandexSearchAlbums: map[string][]byte{},
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
