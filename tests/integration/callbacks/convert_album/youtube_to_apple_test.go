package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumYoutubeToApple(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when youtube album link given and apple album found",
			input: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ap",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://music.apple.com/us/album/amnesiac/1097864180",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				AppleSearchAlbums: map[string][]byte{
					"Radiohead Amnesiac": fixture.Read("apple/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:  "when youtube album link given and apple album not found",
			input: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ap",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Album not found in Apple",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				AppleSearchAlbums: map[string][]byte{},
			},
		},
		{
			name:             "when youtube album not found",
			input:            "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/ap",
			expectedMessages: []*telegram.Message{},
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

			require.Equal(t, tt.expectedMessages, senderMock.AllSent)
		})
	}
}
