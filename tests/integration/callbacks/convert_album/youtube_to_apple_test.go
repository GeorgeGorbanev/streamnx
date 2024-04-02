package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumYoutubeToApple(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedText string
		fixturesMap  fixture.FixturesMap
	}{
		{
			name:         "when youtube album link given and apple album found",
			input:        "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ap",
			expectedText: "https://music.apple.com/us/album/amnesiac/1097864180",
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				// TODO: fix adapter query
				AppleSearchAlbums: map[string][]byte{
					"Radiohead Amnesiac (2001) ": fixture.Read("apple/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:         "when youtube album link given and apple album not found",
			input:        "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ap",
			expectedText: "Album not found in Apple",
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				AppleSearchAlbums: map[string][]byte{},
			},
		},
		{
			name:         "when youtube album not found",
			input:        "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ap",
			expectedText: "",
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums:     map[string][]byte{},
				AppleSearchAlbums: map[string][]byte{},
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
