package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumYandexToSpotify(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  fixture.FixturesMap
	}{
		{
			name:         "when yandex album link given and spotify album found",
			input:        "cnval/ya/3389008/sf",
			expectedText: "https://open.spotify.com/album/1HrMmB5useeZ0F5lHrMvl0",
			fixturesMap: fixture.FixturesMap{
				YandexAlbums: map[string][]byte{
					"3389008": fixture.Read("yandex/get_album_radiohead_amnesiac.json"),
				},
				SpotifySearchAlbums: map[string][]byte{
					"artist:Radiohead album:Amnesiac": fixture.Read("spotify/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:         "when yandex album link given and spotify album not found",
			input:        "cnval/ya/3389008/sf",
			expectedText: "Album not found in Spotify",
			fixturesMap: fixture.FixturesMap{
				YandexAlbums: map[string][]byte{
					"3389008": fixture.Read("yandex/get_album_radiohead_amnesiac.json"),
				},
				SpotifySearchAlbums: map[string][]byte{},
			},
		},
		{
			name:         "when yandex album not found",
			input:        "cnval/ya/3389008/sf",
			expectedText: "",
			fixturesMap: fixture.FixturesMap{
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
