package convert_album

import (
	"testing"

	"github.com/GeorgeGorbanev/vibeshare/internal/telegram"
	"github.com/GeorgeGorbanev/vibeshare/tests/fixture"

	"github.com/stretchr/testify/require"
	"github.com/tucnak/telebot"
)

func TestCallback_ConvertAlbumYoutubeToSpotify(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedMessages []*telegram.Message
		fixturesMap      fixture.FixturesMap
	}{
		{
			name:  "when youtube album link given and spotify album found",
			input: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/sf",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "https://open.spotify.com/album/1HrMmB5useeZ0F5lHrMvl0",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				SpotifySearchAlbums: map[string][]byte{
					"artist:Radiohead album:Amnesiac": fixture.Read("spotify/search_album_radiohead_amnesiac.json"),
				},
			},
		},
		{
			name:  "when youtube album link given and spotify album not found",
			input: "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/sf",
			expectedMessages: []*telegram.Message{
				{
					To:   user,
					Text: "Album not found in Spotify",
				},
			},
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums: map[string][]byte{
					"PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C": fixture.Read("youtube/get_album_radiohead_amnesiac.json"),
				},
				SpotifySearchAlbums: map[string][]byte{},
			},
		},
		{
			name:             "when youtube album not found",
			input:            "cnval/yt/PLAV7kVdctKCbILB72QeXGTVe9DhgnsL0C/sf",
			expectedMessages: []*telegram.Message{},
			fixturesMap: fixture.FixturesMap{
				YoutubeAlbums:       map[string][]byte{},
				SpotifySearchAlbums: map[string][]byte{},
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
