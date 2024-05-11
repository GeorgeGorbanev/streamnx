package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumAppleToSpotify(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when apple album link given and spotify album found",
			input: "cnval/ap/1097864180/sf",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://open.spotify.com/album/1HrMmB5useeZ0F5lHrMvl0",
				},
			},
			fixturesMap: fixture.FixturesMap{
				AppleAlbums: map[string][]byte{
					"1097864180": fixture.Read("apple/get_album_radiohead_amnesiac.json"),
				},
				SpotifySearchAlbums: map[string][]byte{
					"artist:Radiohead album:Amnesiac": fixture.Read("spotify/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:  "when apple album link given and spotify album not found",
			input: "cnval/ap/1097864180/sf",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Album not found in Spotify",
				},
			},
			fixturesMap: fixture.FixturesMap{
				AppleAlbums: map[string][]byte{
					"1097864180": fixture.Read("apple/get_album_radiohead_amnesiac.json"),
				},
				SpotifySearchAlbums: map[string][]byte{},
			},
		},
		{
			name:             "when apple album not found",
			input:            "cnval/ap/1097864180/sf",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				AppleAlbums:        map[string][]byte{},
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

			require.Equal(t, tt.expectedMessages, senderMock.AllSent)
		})
	}
}
